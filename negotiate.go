package webutil

import (
	"strings"

	"github.com/markusthoemmes/goautoneg"
)

// NegotiateContentType returns the best offered content type for the request's
// Accept header. If two offers match with equal weight, then the more specific
// offer is preferred.  For example, text/* trumps */*. If two offers match
// with equal weight and specificity, then the offer earlier in the list is
// preferred. If no offers match, then defaultOffer is returned.
func NegotiateContentType(acceptHeader string, offeredContentTypes []string, defaultOffer string) string {
	acceptedContentTypes := goautoneg.ParseAccept(acceptHeader)

	bestOffer := defaultOffer
	bestQ := -1.0
	bestWild := 3
	for _, offer := range offeredContentTypes {
		for _, spec := range acceptedContentTypes {
			switch {
			case spec.Q == 0.0:
				// ignore
			case spec.Q < bestQ:
				// better match found
			case spec.Type == "*" && spec.SubType == "*":
				if spec.Q > bestQ || bestWild > 2 {
					bestQ = spec.Q
					bestWild = 2
					bestOffer = offer
				}
			case spec.SubType == "*":
				if strings.HasPrefix(offer, spec.Type+"/") && (spec.Q > bestQ || bestWild > 1) {
					bestQ = spec.Q
					bestWild = 1
					bestOffer = offer
				}
			default:
				if spec.Type+"/"+spec.SubType == offer && (spec.Q > bestQ || bestWild > 0) {
					bestQ = spec.Q
					bestWild = 0
					bestOffer = offer
				}
			}
		}
	}
	return bestOffer
}
