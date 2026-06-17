package email

type MockEmailService struct {
	SendEmailFunc func(email string) error
}

func (m *MockEmailService) SendEmail(email string) error {
	if m.SendEmailFunc != nil {
		return m.SendEmailFunc(email)
	}

	return nil
}
