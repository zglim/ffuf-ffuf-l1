package filter

import (
	"fmt"
	"sync"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// MatcherManager handles both filters and matchers.
type MatcherManager struct {
	IsCalibrated     bool
	Mutex            sync.Mutex
	Matchers         map[string]ffuf.FilterProvider
	Filters          map[string]ffuf.FilterProvider
	PerDomainFilters map[string]*PerDomainFilter
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

// AddFilter adds or appends a filter to the PerDomainFilter
func (p *PerDomainFilter) AddFilter(name string, option string) error {
	return addOrUpdate(p.Filters, name, option, false)
}

// GetFilters returns the filters map of the PerDomainFilter
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

func (f *MatcherManager) SetCalibratedForHost(host string, value bool) {
	if f.PerDomainFilters[host] != nil {
		f.PerDomainFilters[host].IsCalibrated = value
	} else {
		newFilter := NewPerDomainFilter(f.Filters)
		newFilter.IsCalibrated = true
		f.PerDomainFilters[host] = newFilter
	}
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

// addOrUpdate creates a new filter and either sets or appends it in the given map.
func addOrUpdate(filters map[string]ffuf.FilterProvider, name string, option string, replace bool) error {
	newf, err := NewFilterByName(name, option)
	if err != nil {
		return err
	}
	if filters[name] == nil || replace {
		filters[name] = newf
	} else {
		newoption := filters[name].Repr() + "," + option
		newerf, err := NewFilterByName(name, newoption)
		if err == nil {
			filters[name] = newerf
		}
	}
	return nil
}

//AddFilter adds a new filter to MatcherManager
func (f *MatcherManager) AddFilter(name string, option string, replace bool) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	return addOrUpdate(f.Filters, name, option, replace)
}

//AddPerDomainFilter adds a new filter to PerDomainFilter configuration
func (f *MatcherManager) AddPerDomainFilter(domain string, name string, option string) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	if _, ok := f.PerDomainFilters[domain]; !ok {
		f.PerDomainFilters[domain] = NewPerDomainFilter(f.Filters)
	}
	return f.PerDomainFilters[domain].AddFilter(name, option)
}

//RemoveFilter removes a filter of a given type
func (f *MatcherManager) RemoveFilter(name string) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	delete(f.Filters, name)
}

//AddMatcher adds a new matcher to Config
func (f *MatcherManager) AddMatcher(name string, option string) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	return addOrUpdate(f.Matchers, name, option, false)
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
	return f.PerDomainFilters[domain].GetFilters()
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

// Match evaluates matchers and filters against a response and returns true if
// the response should be considered a match.
func (f *MatcherManager) Match(resp ffuf.Response, matcherMode string, filterMode string, autoCalibrationPerHost bool) bool {
	matched := false
	for _, m := range f.Matchers {
		match, err := m.Filter(&resp)
		if err != nil {
			continue
		}
		if match {
			matched = true
		} else if matcherMode == "and" {
			return false
		}
	}
	if !matched {
		return false
	}

	var filters map[string]ffuf.FilterProvider
	if autoCalibrationPerHost {
		filters = f.FiltersForDomain(ffuf.HostURLFromRequest(*resp.Request))
	} else {
		filters = f.Filters
	}

	for _, flt := range filters {
		fv, err := flt.Filter(&resp)
		if err != nil {
			continue
		}
		if fv {
			if filterMode == "or" {
				return false
			}
		} else {
			if filterMode == "and" {
				return true
			}
		}
	}
	if len(filters) > 0 && filterMode == "and" {
		return false
	}
	return true
}
