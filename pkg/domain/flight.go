package domain

import "github.com/pkg/errors"

type Flights []*Flight

type Flight struct {
	Source      Airport
	Destination Airport
}

func (f Flights) OriginalSourceAndDestination() (*Flight, error) {
	if len(f) == 0 {
		return nil, ErrEmptyFlightsList
	}

	var (
		initialSource    Airport
		finalDestination Airport
		allSources       = make(map[Airport]struct{})
		allDestinations  = make(map[Airport]struct{})
	)

	for _, v := range f {
		if _, ok := allSources[v.Source]; ok {
			return nil, errors.Wrapf(ErrInvalidItinerary, "'%v' source appears more than once in the sources", v.Source)
		}

		if _, ok := allDestinations[v.Destination]; ok {
			return nil, errors.Wrapf(ErrInvalidItinerary, "'%v' destination appears more than once in the destinations", v.Source)
		}

		allSources[v.Source] = struct{}{}
		allDestinations[v.Destination] = struct{}{}
	}

	for source := range allSources {
		if _, ok := allDestinations[source]; !ok {
			initialSource = source
		}
	}

	for destination := range allDestinations {
		if _, ok := allSources[destination]; !ok {
			finalDestination = destination
		}
	}

	return &Flight{
		Source:      initialSource,
		Destination: finalDestination,
	}, nil
}
