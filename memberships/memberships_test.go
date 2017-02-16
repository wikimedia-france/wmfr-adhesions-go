package memberships

import (
	"testing"
	"github.com/wikimedia-france/wmfr-adhesions/iraiser"
	"time"
	"math/rand"
	"github.com/stretchr/testify/assert"
	"github.com/wikimedia-france/wmfr-adhesions/civicrm"
	"github.com/wikimedia-france/wmfr-adhesions/internal"
)

func TestRecordMembershipRenewal(t *testing.T) {
	const contactMail = "test@example.org"
	contactId := rand.Int()
	membershipId := rand.Int()
	contributionId := rand.Int()
	membershipTypeId := rand.Int()

	validationDate := time.Now()
	amount := rand.Int()

	donation := iraiser.Donation{
		Donator: iraiser.Donator{
			Mail: contactMail,
			FirstName: "First",
			LastName: "Last",
			Pseudo: "Nick",
			StreetAddress: "Address",
			City: "City",
			PostalCode: "12345",
			Country: "CY",
		},
		Campaign: iraiser.Campaign{
			AffectationCode: "Aff",
			OriginCode: "Ori",
		},
		Payment: iraiser.Payment{
			Mode: "card",
			GatewayId: "6789",
		},
		Amount: amount,
		Currency: "CUR",
		Reference: "abcd",
		ValidationDate: validationDate,
	}
	internal.Config.MembershipTypeId = membershipTypeId

	contactGetter = func(query *civicrm.GetContactQuery) (*civicrm.GetContactResponse, error) {
		assert.Equal(t, contactMail, query.Mail)
		return &civicrm.GetContactResponse{
			StatusResponse: civicrm.StatusResponse{
				Id: contactId,
				Count: 1,
			},
		}, nil
	}

	membershipGetter = func(query *civicrm.GetMembershipQuery) (*civicrm.GetMembershipResponse, error) {
		assert.Equal(t, contactId, query.ContactId)
		return &civicrm.GetMembershipResponse{
			Values: map[int]civicrm.Membership {
				0: {
					Id: membershipId,
					ContactId: contactId,
					MembershipTypeId: membershipTypeId,
				},
			},
		}, nil
	}

	membershipCreator = func(membership *civicrm.Membership) (*civicrm.CreateMembershipResponse, error) {
		assert.Equal(t, membershipId, membership.Id)
		assert.Equal(t, contactId, membership.ContactId)
		assert.Equal(t, StatusOverride, membership.StatusOverride)
		assert.Equal(t, Terms, membership.Terms)
		return &civicrm.CreateMembershipResponse{}, nil
	}

	contributionCreator = func(contribution *civicrm.Contribution) (*civicrm.CreateContributionResponse, error) {
		assert.Equal(t, contactId, contribution.ContactId)
		return &civicrm.CreateContributionResponse{
			StatusResponse: civicrm.StatusResponse{
				Id: contributionId,
			},
		}, nil
	}

	membershipPaymentCreator = func(payment *civicrm.MembershipPayment) (*civicrm.CreateMembershipPaymentResponse, error) {
		assert.Equal(t, contributionId, payment.ContributionId)
		assert.Equal(t, membershipId, payment.MembershipId)
		return &civicrm.CreateMembershipPaymentResponse{}, nil
	}

	_, err := RecordMembership(&donation)
	assert.NoError(t, err)
}
