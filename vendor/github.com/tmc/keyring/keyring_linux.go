// +build !gnome_keyring

package keyring

import (
	"fmt"

	dbus "github.com/guelfey/go.dbus"
)

const (
	ssServiceName     = "org.freedesktop.secrets"
	ssServicePath     = "/org/freedesktop/secrets"
	ssCollectionPath  = "/org/freedesktop/secrets/collection/Default"
	ssServiceIface    = "org.freedesktop.Secret.Service."
	ssSessionIface    = "org.freedesktop.Secret.Session."
	ssCollectionIface = "org.freedesktop.Secret.Collection."
	ssItemIface       = "org.freedesktop.Secret.Item."
	ssPromptIface     = "org.freedesktop.Secret.Prompt."
)

// ssSecret corresponds to org.freedesktop.Secret.Item
// Note: Order is important
type ssSecret struct {
	Session     dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string `dbus:"content_type"`
}

// newSSSecret prepares an ssSecret for use
// Uses text/plain as the Content-type which may need to change in the future
func newSSSecret(session dbus.ObjectPath, secret string) (s ssSecret) {
	s = ssSecret{
		ContentType: "text/plain; charset=utf8",
		Parameters:  []byte{},
		Session:     session,
		Value:       []byte(secret),
	}
	return
}

// ssProvider implements the provider interface freedesktop SecretService
type ssProvider struct {
	*dbus.Conn
	srv *dbus.Object
}

// This is used to open a seassion for every get/set. Alternative might be to
// defer() the call to close when constructing the ssProvider
func (s *ssProvider) openSession() (*dbus.Object, error) {
	var disregard dbus.Variant
	var sessionPath dbus.ObjectPath
	method := fmt.Sprint(ssServiceIface, "OpenSession")
	err := s.srv.Call(method, 0, "plain", dbus.MakeVariant("")).Store(&disregard, &sessionPath)
	if err != nil {
		return nil, err
	}
	return s.Object(ssServiceName, sessionPath), nil
}

// Unsure how the .Prompt call surfaces, it hasn't come up.
func (s *ssProvider) unlock(p dbus.ObjectPath) error {
	var unlocked []dbus.ObjectPath
	var prompt dbus.ObjectPath
	method := fmt.Sprint(ssServiceIface, "Unlock")
	err := s.srv.Call(method, 0, []dbus.ObjectPath{p}).Store(&unlocked, &prompt)
	if err != nil {
		return fmt.Errorf("keyring/dbus: Unlock error: %s", err)
	}
	if prompt != dbus.ObjectPath("/") {
		method = fmt.Sprint(ssPromptIface, "Prompt")
		call := s.Object(ssServiceName, prompt).Call(method, 0, "unlock")
		return call.Err
	}
	return nil
}

func (s *ssProvider) Get(c, u string) (string, error) {
	results := []dbus.ObjectPath{}
	var secret ssSecret
	search := map[string]string{
		"username": u,
		"service":  c,
	}

	session, err := s.openSession()
	if err != nil {
		return "", err
	}
	defer session.Call(fmt.Sprint(ssSessionIface, "Close"), 0)
	s.unlock(ssCollectionPath)
	collection := s.Object(ssServiceName, ssCollectionPath)

	method := fmt.Sprint(ssCollectionIface, "SearchItems")
	call := collection.Call(method, 0, search)
	err = call.Store(&results)
	if call.Err != nil {
		return "", call.Err
	}
	// results is a slice. Just grab the first one.
	if len(results) == 0 {
		return "", ErrNotFound
	}

	method = fmt.Sprint(ssItemIface, "GetSecret")
	err = s.Object(ssServiceName, results[0]).Call(method, 0, session.Path()).Store(&secret)
	if err != nil {
		return "", err
	}
	return string(secret.Value), nil
}

func (s *ssProvider) Set(c, u, p string) error {
	var item, prompt dbus.ObjectPath
	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label": dbus.MakeVariant(fmt.Sprintf("%s - %s", u, c)),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
			"username": u,
			"service":  c,
		}),
	}

	session, err := s.openSession()
	if err != nil {
		return err
	}
	defer session.Call(fmt.Sprint(ssSessionIface, "Close"), 0)
	s.unlock(ssCollectionPath)
	collection := s.Object(ssServiceName, ssCollectionPath)

	secret := newSSSecret(session.Path(), p)
	// the bool is "replace"
	err = collection.Call(fmt.Sprint(ssCollectionIface, "CreateItem"), 0, properties, secret, true).Store(&item, &prompt)
	if err != nil {
		return fmt.Errorf("keyring/dbus: CreateItem error: %s", err)
	}
	if prompt != "/" {
		s.Object(ssServiceName, prompt).Call(fmt.Sprint(ssPromptIface, "Prompt"), 0, "unlock")
	}
	return nil
}

func initializeProvider() (provider, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	srv := conn.Object(ssServiceName, ssServicePath)
	p := &ssProvider{conn, srv}
	return p, nil
}
