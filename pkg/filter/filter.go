package filter

import (
	"fmt"
	"sync"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// MatcherManager handles both filters and matchers.
type MatcherManager struct {
	IsCalibrated          bool
	Mutex                 sync.Mutex
	Matchers              map[string]ffuf.FilterProvider
	Filters               map[string]ffuf.FilterProvider
	PerDomainFilters      map[string]*PerDomainFilter
	MatcherMode           string
	FilterMode            string
	AutoCalibrationPerHost bool
}

type PerDomainFilter struct {
	IsCalibrated bool
	Filters      map[string]ffuf.FilterProvider
}

func NewPerDomainFilter(globfilters map[string]ffuf.FilterProvider) *PerDomainFilter {
	return &PerDomainFilter{IsCalibrated: false, Filters: globfilters}
}

func (p *PerDomainFilter) SetCalibrated(value bool) {
	p.IsCalibrated = value
}

// AddFilter adds a filter to the PerDomainFilter, creating or appending as needed.
func (p *PerDomainFilter) AddFilter(name string, option string) error {
	newf, err := NewFilterByName(name, option)
	if err != nil {
		return err
	}
	if p.Filters[name] == nil {
		p.Filters[name] = newf
	} else {
		newoption := p.Filters[name].Repr() + "," + option
		newerf, err := NewFilterByName(name, newoption)
		if err != nil {
			return err
		}
		p.Filters[name] = newerf
	}
	return nil
}

// GetFilters returns the filters for this PerDomainFilter.
func (p *PerDomainFilter) GetFilters() map[string]ffuf.FilterProvider {
	return p.Filters
}

func NewMatcherManager() ffuf.MatcherManager {
	return &MatcherManager{
		IsCalibrated:     false,
		Matchers:         make(map[string]ffuf.FilterProvider),
		Filters:          make(map[string]ffuf.FilterProvider),
		PerDomainFilters: make(map[string]*PerDomainFilter),
	}
}

func (f *MatcherManager) SetCalibrated(value bool) {
	f.IsCalibrated = value
}

// SetMatcherMode sets the matcher mode ("and" or "or").
func (f *MatcherManager) SetMatcherMode(mode string) {
	f.MatcherMode = mode
}

// SetFilterMode sets the filter mode ("and" or "or").
func (f *MatcherManager) SetFilterMode(mode string) {
	f.FilterMode = mode
}

// SetAutoCalibrationPerHost sets whether per-host autocalibration is active.
func (f *MatcherManager) SetAutoCalibrationPerHost(value bool) {
	f.AutoCalibrationPerHost = value
}

// getOrCreatePerDomainFilter returns the PerDomainFilter for the given domain,
// creating and storing a new one (initialized with global filters) if none exists.
func (f *MatcherManager) getOrCreatePerDomainFilter(domain string) *PerDomainFilter {
	if pd, ok := f.PerDomainFilters[domain]; ok {
		return pd
	}
	pd := NewPerDomainFilter(f.Filters)
	f.PerDomainFilters[domain] = pd
	return pd
}

func (f *MatcherManager) SetCalibratedForHost(host string, value bool) {
	pd := f.getOrCreatePerDomainFilter(host)
	pd.IsCalibrated = value
}

func NewFilterByName(name string, value string) (ffuf.FilterProvider, error) {
	if name == "status" {
		return NewStatusFilter(value)
	}
	if name == "size" {
		return NewSizeFilter(value)
	}
	if name == "word" {
		return NewWordFilter(value)
	}
	if name == "line" {
		return NewLineFilter(value)
	}
	if name == "regexp" {
		return NewRegexpFilter(value)
	}
	if name == "time" {
		return NewTimeFilter(value)
	}
	return nil, fmt.Errorf("Could not create filter with name %s", name)
}

// addOrUpdate is a private helper that encapsulates the shared
// "create → check existence → append or replace" logic used by
// both AddFilter and AddMatcher.
func (f *MatcherManager) addOrUpdate(collection map[string]ffuf.FilterProvider, name string, option string, replace bool) error {
	newf, err := NewFilterByName(name, option)
	if err != nil {
		return err
	}
	if collection[name] == nil || replace {
		collection[name] = newf
	} else {
		newoption := collection[name].Repr() + "," + option
		newerf, err := NewFilterByName(name, newoption)
		if err != nil {
			return err
		}
		collection[name] = newerf
	}
	return nil
}

// AddFilter adds a new filter to MatcherManager
func (f *MatcherManager) AddFilter(name string, option string, replace bool) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	return f.addOrUpdate(f.Filters, name, option, replace)
}

// AddPerDomainFilter adds a new filter to PerDomainFilter configuration
func (f *MatcherManager) AddPerDomainFilter(domain string, name string, option string) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	pd := f.getOrCreatePerDomainFilter(domain)
	return pd.AddFilter(name, option)
}

// RemoveFilter removes a filter of a given type
func (f *MatcherManager) RemoveFilter(name string) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	delete(f.Filters, name)
}

// AddMatcher adds a new matcher to Config
func (f *MatcherManager) AddMatcher(name string, option string) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	return f.addOrUpdate(f.Matchers, name, option, false)
}

func (f *MatcherManager) GetFilters() map[string]ffuf.FilterProvider {
	return f.Filters
}

func (f *MatcherManager) GetMatchers() map[string]ffuf.FilterProvider {
	return f.Matchers
}

func (f *MatcherManager) FiltersForDomain(domain string) map[string]ffuf.FilterProvider {
	if f.PerDomainFilters[domain] == nil {
		return f.Filters
	}
	return f.PerDomainFilters[domain].Filters
}

func (f *MatcherManager) CalibratedForDomain(domain string) bool {
	if f.PerDomainFilters[domain] != nil {
		return f.PerDomainFilters[domain].IsCalibrated
	}
	return false
}

func (f *MatcherManager) Calibrated() bool {
	return f.IsCalibrated
}

// Match evaluates whether a response passes the matcher/filter logic.
// It encapsulates the and/or branching for both matchers and filters,
// as well as the AutoCalibrationPerHost domain-filter selection.
func (f *MatcherManager) Match(resp ffuf.Response) bool {
	// Determine which filter set to use
	var filters map[string]ffuf.FilterProvider
	if f.AutoCalibrationPerHost && resp.Request != nil {
		filters = f.FiltersForDomain(ffuf.HostURLFromRequest(*resp.Request))
	} else {
		filters = f.Filters
	}
	matchers := f.Matchers

	// Evaluate matchers
	matched := false
	for _, m := range matchers {
		match, err := m.Filter(&resp)
		if err != nil {
			continue
		}
		if match {
			matched = true
		} else if f.MatcherMode == "and" {
			// In "and" mode a single non-match means the whole set fails
			return false
		}
	}
	// If no matcher matched, the response is not a match
	if !matched {
		return false
	}

	// Evaluate filters
	for _, flt := range filters {
		fv, err := flt.Filter(&resp)
		if err != nil {
			continue
		}
		if fv {
			if f.FilterMode == "or" {
				// In "or" mode a single filter match excludes the response
				return false
			}
		} else {
			if f.FilterMode == "and" {
				// In "and" mode not all filters matched, response passes
				return true
			}
		}
	}
	if len(filters) > 0 && f.FilterMode == "and" {
		// All filters matched in "and" mode, exclude the response
		return false
	}
	return true
}
