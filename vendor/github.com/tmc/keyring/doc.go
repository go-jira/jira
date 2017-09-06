/*
Package keyring provides a cross-platform interface to keychains for
password management

Currently implemented:

 * OSX
 * SecretService
 * gnome-keychain (via "gnome_keyring" build flag)


Usage

Example usage:


   err := keyring.Set("libraryFoo", "jack", "sacrifice")
   password, err := keyring.Get("libraryFoo", "jack")
   fmt.Println(password)
   Output: sacrifice


TODO

    * Write Windows provider
*/
package keyring
