package memberships

import (
	"github.com/wikimedia-france/wmfr-adhesions/iraiser"
	"github.com/wikimedia-france/wmfr-adhesions/civicrm"
)

func (m *Memberships) getContact(donation *iraiser.Donation) (int, error) {
	query := civicrm.GetContactQuery{
		Mail: donation.Donator.Mail,
	}
	searchResult, err := m.crm.GetContact(&query);
	if err != nil {
		return -1, err
	}
	if searchResult.Count == 0 {
		return -1, &NoSuchContactError{Mail: donation.Donator.Mail}
	}
	if searchResult.Count > 1 {
		return -1, &TooManyContactsError{Mail: donation.Donator.Mail, Count: searchResult.Count}
	}
	return searchResult.Id, nil
}

func (m *Memberships) createContact(donation *iraiser.Donation) (int, error) {
	contact := civicrm.Contact{
		ContactType: m.config.ContactTypeName,
		Mail: donation.Donator.Mail,
		FirstName: donation.Donator.FirstName,
		LastName: donation.Donator.LastName,
		Pseudo: donation.Donator.Pseudo,
		Source: m.config.ContactSourceName,
	}
	contactResp, err := m.crm.CreateContact(&contact)
	if err != nil {
		return -1, err
	}
	address := civicrm.Address{
		ContactId: contactResp.Id,
		LocationTypeId: m.config.LocationTypeId,
		StreetAddress: donation.Donator.StreetAddress,
		City: donation.Donator.City,
		PostalCode: donation.Donator.PostalCode,
		Country: donation.Donator.Country,
		StreetParsing: 1,
	}
	_, err = m.crm.CreateAddress(&address)
	return contactResp.Id, err
}
