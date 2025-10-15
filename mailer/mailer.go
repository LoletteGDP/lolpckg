// mailer/mailer.go
package mailer

import (
	gomail "gopkg.in/mail.v2"
)

// Mailer encapsula el Dialer de gomail y la dirección del remitente.
type Mailer struct {
	dialer *gomail.Dialer
	sender string // Dirección de correo desde la que se enviarán los emails (ej: "tu@empresa.com")
}

// New crea una nueva instancia de Mailer.
// Necesita el host, puerto, usuario y contraseña de tu servidor SMTP.
func New(host string, port int, username, password, sender string) Mailer {
	// gomail.NewDialer crea un "marcador" que se puede reutilizar para conectar
	// al servidor SMTP.
	d := gomail.NewDialer(host, port, username, password)

	// Guardamos el marcador y la dirección del remitente para usarlos después.
	return Mailer{
		dialer: d,
		sender: sender,
	}
}

// Send es el método para enviar un correo.
// Recibe el destinatario, el asunto y el cuerpo del mensaje en formato HTML.
func (m Mailer) Send(to, subject, htmlBody string) error {
	// Creamos un nuevo mensaje.
	msg := gomail.NewMessage()

	// Establecemos los encabezados del correo.
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	// Establecemos el cuerpo del mensaje. Usamos "text/html" para poder enviar
	// correos con formato, enlaces, imágenes, etc.
	msg.SetBody("text/html", htmlBody)

	// Finalmente, usamos el Dialer para conectar con el servidor SMTP y enviar el mensaje.
	// DialAndSend se encarga de abrir la conexión, enviar y cerrarla.
	return m.dialer.DialAndSend(msg)
}
