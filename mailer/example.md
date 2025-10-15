package main

import (
	"fmt"
	"log"

	// Importa tu nuevo paquete mailer
	"github.com/LoletteGDP/lolpckg/mailer"
)

func main() {
	// --- Configuración del Mailer ---
	// ¡IMPORTANTE! Nunca pongas credenciales directamente en el código.
	// Lo ideal es leerlas desde variables de entorno o un archivo de configuración.
	smtpHost := "smtp.example.com" // El host de tu proveedor (ej: smtp.gmail.com)
	smtpPort := 587                // El puerto (normalmente 587 para TLS)
	smtpUser := "tu_usuario"       // Tu email o usuario SMTP
	smtpPass := "tu_contraseña"    // Tu contraseña o clave de aplicación
	senderMail := "no-reply@tuapp.com"

	// 1. Creas una instancia de tu Mailer
	myMailer := mailer.New(smtpHost, smtpPort, smtpUser, smtpPass, senderMail)

	// 2. Preparas los datos del correo a enviar
	recipient := "usuario.nuevo@email.com"
	subject := "¡Bienvenido a nuestra aplicación!"
	body := "<h1>Hola y bienvenido!</h1><p>Estamos muy contentos de tenerte con nosotros.</p>"

	// 3. Envías el correo usando el método Send
	err := myMailer.Send(recipient, subject, body)
	if err != nil {
		log.Fatalf("No se pudo enviar el correo: %s", err)
	}

	fmt.Println("¡Correo enviado exitosamente!")
}