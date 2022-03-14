package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	mockpolicy "github.com/MagalixCorp/magalix-policy-agent/internal/policies/mock"
	mocksink "github.com/MagalixCorp/magalix-policy-agent/internal/sink/mock"
	"github.com/MagalixCorp/magalix-policy-agent/pkg/domain"
	"github.com/MagalixCorp/magalix-policy-agent/pkg/validation/testdata"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewOPAValidator(t *testing.T) {
	type args struct {
		policiesSource  domain.PoliciesSource
		writeCompliance bool
		resultsSinks    []domain.PolicyValidationSink
		validationType  string
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	policiesSource := mockpolicy.NewMockPoliciesSource(ctrl)
	sink := mocksink.NewMockPolicyValidationSink(ctrl)
	tests := []struct {
		name string
		args args
		want *OpaValidator
	}{
		{
			name: "default test",
			args: args{
				policiesSource:  policiesSource,
				writeCompliance: true,
				resultsSinks:    []domain.PolicyValidationSink{sink},
				validationType:  "TestValidate",
			},
			want: &OpaValidator{
				policiesSource:  policiesSource,
				writeCompliance: true,
				resultsSinks:    []domain.PolicyValidationSink{sink},
				validationType:  "TestValidate",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOPAValidator(tt.args.policiesSource, tt.args.writeCompliance, tt.args.validationType, tt.args.resultsSinks...); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NewOPAValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func cmpPolicyValidation(arg1, arg2 domain.PolicyValidation) bool {
	if arg1.Entity.ID != arg2.Entity.ID {
		return false
	}

	return arg1.Type == arg2.Type && arg1.Trigger == arg2.Trigger && arg1.Status == arg2.Status
}

func getEntityFromStringSpec(entityStringSpec string) (domain.Entity, error) {
	var entitySpec map[string]interface{}
	err := json.Unmarshal([]byte(entityStringSpec), &entitySpec)
	if err != nil {
		return domain.Entity{}, fmt.Errorf("invalid string format: %w", err)
	}
	return domain.NewEntityFromSpec(entitySpec), nil
}

func TestOpaValidator_Validate(t *testing.T) {
	type init struct {
		loadStubs       func(*mockpolicy.MockPoliciesSource, *mocksink.MockPolicyValidationSink)
		writeCompliance bool
	}
	assert := require.New(t)

	entityText := testdata.Entity
	validationType := "unit-test"
	entity, err := getEntityFromStringSpec(entityText)
	assert.Nil(err)
	compliantEntity, err := getEntityFromStringSpec(testdata.CompliantEntity)
	assert.Nil(err)
	tests := []struct {
		name    string
		init    init
		entity  domain.Entity
		want    *domain.PolicyValidationSummary
		wantErr bool
	}{
		{
			name: "default test",
			init: init{
				writeCompliance: false,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						testdata.Policies["imageTag"],
						testdata.Policies["missingOwner"],
					}, nil)
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(1).Return(nil)
				},
			},
			entity: entity,
			want: &domain.PolicyValidationSummary{
				Violations: []domain.PolicyValidation{
					{
						Policy: testdata.Policies["imageTag"],
						Entity: entity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusViolating,
					},
					{
						Policy: testdata.Policies["missingOwner"],
						Entity: entity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusViolating,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error getting policies",
			init: init{
				writeCompliance: false,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return(nil, fmt.Errorf(""))
					// expect 0 calls to sink write
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(0).Return(nil)
				},
			},
			entity:  entity,
			want:    nil,
			wantErr: true,
		},
		{
			name: "entity kind matching",
			init: init{
				writeCompliance: false,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					missingOwner := testdata.Policies["missingOwner"]
					missingOwner.Targets = domain.PolicyTargets{Kind: []string{"Deployment"}}
					imageTag := testdata.Policies["imageTag"]
					imageTag.Targets = domain.PolicyTargets{Kind: []string{"ReplicaSet"}}
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						missingOwner,
						imageTag,
					}, nil)
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(1).Return(nil)
				},
			},
			entity: entity,
			want: &domain.PolicyValidationSummary{
				Violations: []domain.PolicyValidation{
					{
						Policy: testdata.Policies["missingOwner"],
						Entity: entity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusViolating,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "entity namespace matching",
			init: init{
				writeCompliance: false,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					missingOwner := testdata.Policies["missingOwner"]
					missingOwner.Targets = domain.PolicyTargets{Namespace: []string{"unit-testing"}}
					imageTag := testdata.Policies["imageTag"]
					imageTag.Targets = domain.PolicyTargets{Namespace: []string{"bad-namespace"}}
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						missingOwner,
						imageTag,
					}, nil)
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(1).Return(nil)
				},
			},
			entity: entity,
			want: &domain.PolicyValidationSummary{
				Violations: []domain.PolicyValidation{
					{
						Policy: testdata.Policies["missingOwner"],
						Entity: entity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusViolating,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "entity labels matching",
			init: init{
				writeCompliance: false,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					missingOwner := testdata.Policies["missingOwner"]
					missingOwner.Targets = domain.PolicyTargets{Label: []map[string]string{{"app": "nginx"}}}
					imageTag := testdata.Policies["imageTag"]
					imageTag.Targets = domain.PolicyTargets{Label: []map[string]string{{"app": "notfound"}}}
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						missingOwner,
						imageTag,
					}, nil)
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(1).Return(nil)
				},
			},
			entity: entity,
			want: &domain.PolicyValidationSummary{
				Violations: []domain.PolicyValidation{
					{
						Policy: testdata.Policies["missingOwner"],
						Entity: entity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusViolating,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple policies only one matching",
			init: init{
				writeCompliance: false,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					missingOwner := testdata.Policies["missingOwner"]
					missingOwner.Targets = domain.PolicyTargets{
						Label:     []map[string]string{{"app": "nginx"}},
						Namespace: []string{"unit-testing"},
						Kind:      []string{"Deployment"},
					}
					imageTag := testdata.Policies["imageTag"]
					imageTag.Targets = domain.PolicyTargets{
						Label:     []map[string]string{{"app": "nginx"}},
						Namespace: []string{"bad-namespace"},
						Kind:      []string{"Deployment"},
					}
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						missingOwner,
						imageTag,
					}, nil)
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(1).Return(nil)
				},
			},
			entity: entity,
			want: &domain.PolicyValidationSummary{
				Violations: []domain.PolicyValidation{
					{
						Policy: testdata.Policies["missingOwner"],
						Entity: entity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusViolating,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no mathching no writes to sink",
			init: init{
				writeCompliance: true,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					imageTag := testdata.Policies["imageTag"]
					imageTag.Targets = domain.PolicyTargets{Label: []map[string]string{{"app": "notfound"}}}
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						imageTag,
					}, nil)
					// expect 0 calls to sink write, no compliance or violation
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(0).Return(nil)
				},
			},
			entity:  entity,
			want:    &domain.PolicyValidationSummary{},
			wantErr: false,
		},
		{
			name: "compliant entity",
			init: init{
				writeCompliance: true,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						testdata.Policies["imageTag"],
						testdata.Policies["missingOwner"],
					}, nil)
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(1).Return(nil)
				},
			},
			entity: compliantEntity,
			want: &domain.PolicyValidationSummary{
				Compliances: []domain.PolicyValidation{
					{
						Policy: testdata.Policies["imageTag"],
						Entity: compliantEntity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusCompliant,
					},
					{
						Policy: testdata.Policies["missingOwner"],
						Entity: compliantEntity,
						Type:   validationType,
						Status: domain.PolicyValidationStatusCompliant,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error loading policy code",
			init: init{
				writeCompliance: false,
				loadStubs: func(policiesSource *mockpolicy.MockPoliciesSource, sink *mocksink.MockPolicyValidationSink) {
					policiesSource.EXPECT().GetAll(gomock.Any()).
						Times(1).Return([]domain.Policy{
						testdata.Policies["badPolicyCode"],
					}, nil)
					sink.EXPECT().Write(gomock.Any(), gomock.Any()).
						Times(0).Return(nil)
				},
			},
			entity:  entity,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			policiesSource := mockpolicy.NewMockPoliciesSource(ctrl)
			sink := mocksink.NewMockPolicyValidationSink(ctrl)
			tt.init.loadStubs(policiesSource, sink)
			v := &OpaValidator{
				policiesSource:  policiesSource,
				resultsSinks:    []domain.PolicyValidationSink{sink},
				writeCompliance: tt.init.writeCompliance,
				validationType:  validationType,
			}
			got, err := v.Validate(context.Background(), tt.entity, validationType)
			if tt.wantErr {
				assert.Error(err)
				return
			}
			if err != nil {
				assert.Fail("validator test failed", "got unexcpected error, %s", err)
			}
			assert.Equal(len(got.Violations), len(tt.want.Violations))
			assert.Equal(len(got.Compliances), len(tt.want.Compliances))

			for _, wantViolation := range tt.want.Violations {
				found := false
				for _, gotViolation := range got.Violations {
					if gotViolation.Policy.ID == wantViolation.Policy.ID {
						found = true
						assert.True(
							cmpPolicyValidation(gotViolation, wantViolation),
							"gotten violation not as expected for policy %s",
							wantViolation.Policy.ID,
						)
					}
				}
				assert.True(found, "did not find violation for policy %s", wantViolation.Policy.Name)
			}

			for _, wantCompliance := range tt.want.Compliances {
				found := false
				for _, gotCompliance := range got.Compliances {
					if gotCompliance.Policy.ID == wantCompliance.Policy.ID {
						found = true
						assert.True(
							cmpPolicyValidation(gotCompliance, wantCompliance),
							"gotten compliance not as expected for policy %s",
							wantCompliance.Policy.ID,
						)
					}
				}
				assert.True(found, "did not find compliance for policy %s", wantCompliance.Policy.Name)
			}
		})
	}
}
