package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_GetMode_GetsModeFromHoverfly(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV6{
		v2.DataViewV6{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV6{
				{
					RequestMatcher: v2.RequestMatcherViewV6{
						Method: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/hoverfly/mode",
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 200,
						Body: `{
							"mode": "test-mode",
							"arguments" : {
								"matchingStrategy":"first",
								"headersWhitelist":["foo","bar"],
								"stateful": true
							}
						}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	}, false)

	mode, err := GetMode(target)
	Expect(err).To(BeNil())

	Expect(mode.Mode).To(Equal("test-mode"))
	Expect(*mode.Arguments.MatchingStrategy).To(Equal("first"))
	Expect(mode.Arguments.Headers).To(Equal([]string{"foo", "bar"}))
	Expect(mode.Arguments.Stateful).To(BeTrue())
}

func Test_GetMode_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := GetMode(inaccessibleTarget)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_GetMode_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV6{
		v2.DataViewV6{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV6{
				{
					RequestMatcher: v2.RequestMatcherViewV6{
						Method: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/hoverfly/mode",
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 400,
						Body:   `{"error": "test error"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	}, false)

	_, err := GetMode(target)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not retrieve mode\n\ntest error"))
}

func Test_SetMode_SendsCorrectHTTPRequest(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV6{
		v2.DataViewV6{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV6{
				{
					RequestMatcher: v2.RequestMatcherViewV6{
						Method: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "PUT",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/hoverfly/mode",
							},
						},
						Body: []v2.MatcherViewV6{
							{
								Matcher: matchers.Json,
								Value:   `{"mode":"capture","arguments":{}}`,
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 200,
						Body:   `{"mode": "capture"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	}, false)

	mode, err := SetModeWithArguments(target, &v2.ModeView{
		Mode: "capture",
	})
	Expect(err).To(BeNil())

	Expect(mode).To(Equal("capture"))
}

func Test_SetMode_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := SetModeWithArguments(inaccessibleTarget, &v2.ModeView{
		Mode: "capture",
	})

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_SetMode_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV6{
		v2.DataViewV6{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV6{
				{
					RequestMatcher: v2.RequestMatcherViewV6{
						Method: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "PUT",
							},
						},
						Path: []v2.MatcherViewV6{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/hoverfly/mode",
							},
						},
					},
					Response: v2.ResponseDetailsViewV6{
						Status: 400,
						Body:   `{"error": "test error"}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	}, false)

	_, err := SetModeWithArguments(target, &v2.ModeView{
		Mode: "capture",
	})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not set mode\n\ntest error"))
}
