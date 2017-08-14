// +build gnome_keyring

package keyring

/*
#cgo pkg-config: libsecret-1 glib-2.0
#include <stdlib.h>
#include "libsecret/secret.h"

SecretSchema keyring_schema =
  {
    "org.github.tmc.keyring.Password",
    SECRET_SCHEMA_NONE,
    {
      { "username", SECRET_SCHEMA_ATTRIBUTE_STRING },
      { "service",  SECRET_SCHEMA_ATTRIBUTE_STRING },
      {  NULL, 0 },
    }
  };

// wrap the gnome calls because cgo can't deal with vararg functions

gboolean gkr_set_password(gchar *description, gchar *service, gchar *username, gchar *password, GError **err) {
	return secret_password_store_sync(
		&keyring_schema,
		NULL,
		description,
		password,
    NULL,
    err,
		"service", service,
		"username", username,
		NULL);
}

gchar * gkr_get_password(gchar *service, gchar *username, GError **err) {
	return secret_password_lookup_sync(
		&keyring_schema,
    NULL,
		err,
		"service", service,
		"username", username,
		NULL);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type gnomeKeyring struct{}

func (p gnomeKeyring) Set(Service, Username, Password string) error {
	desc := (*C.gchar)(C.CString("Username and password for " + Service))
	username := (*C.gchar)(C.CString(Username))
	service := (*C.gchar)(C.CString(Service))
	password := (*C.gchar)(C.CString(Password))
	defer C.free(unsafe.Pointer(desc))
	defer C.free(unsafe.Pointer(username))
	defer C.free(unsafe.Pointer(service))
	defer C.free(unsafe.Pointer(password))

	var gerr *C.GError
	result := C.gkr_set_password(desc, service, username, password, &gerr)
	defer C.free(unsafe.Pointer(gerr))

	if result == 0 {
		return fmt.Errorf("Gnome-keyring error: %+v", gerr)
	}
	return nil
}

func (p gnomeKeyring) Get(Service string, Username string) (string, error) {
	var gerr *C.GError
	var pw *C.gchar

	username := (*C.gchar)(C.CString(Username))
	service := (*C.gchar)(C.CString(Service))
	defer C.free(unsafe.Pointer(username))
	defer C.free(unsafe.Pointer(service))

	pw = C.gkr_get_password(service, username, &gerr)
	defer C.free(unsafe.Pointer(gerr))
	defer C.secret_password_free((*C.gchar)(pw))

	if pw == nil {
		return "", fmt.Errorf("Gnome-keyring error: %+v", gerr)
	}
	return C.GoString((*C.char)(pw)), nil
}

func initializeProvider() (provider, error) {
	return gnomeKeyring{}, nil
}
