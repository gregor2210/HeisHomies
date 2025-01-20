package main

import (
	"bufio"  // Lar oss lese tekst, for eksempel det du skriver i terminalen
	"bytes"  // Hjelper oss med å rydde opp i data, fjerne unødvendige ting
	"fmt"    // Brukes til å skrive ut ting på skjermen
	"net"    // Brukes til å lage nettverksforbindelser, som til en server
	"os"     // Lar oss få ting som brukerinndata fra terminalen
)

// Lager forbindelse til serveren
func connectToServer(serverIP string, serverPort int) (net.Conn, error) {
	fmt.Println("Connecting to server...") // Forteller hva vi gjør
	// Prøv å koble til serveren med IP-adresse og port
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		// Hvis det ikke fungerer, gi beskjed om feilen
		return nil, fmt.Errorf("error connecting to server: %v", err)
	}
	fmt.Println("Connected to server") // Forteller at det gikk bra
	return conn, nil
}

// Leser en melding fra serveren
func receiveMessage(conn net.Conn, bufferSize int) (string, error) {
	// Lag en "postkasse" (buffer) som kan holde meldinger på bufferSize (her 1024) bytes
	buffer := make([]byte, bufferSize)
	// Prøv å lese meldingen fra serveren
	n, err := conn.Read(buffer)
	if err != nil {
		// Hvis det er en feil, si hva som gikk galt
		return "", fmt.Errorf("error reading message: %v", err)
	}
	// Trim vekk "usynlige" tegn (nullbytes) fra meldingen
	return string(bytes.Trim(buffer[:n], "\x00")), nil
}

// Sender en melding til serveren
func sendMessage(conn net.Conn, message string, bufferSize int) error {
	// Lag en "postkasse" for meldingen, som alltid fylles opp til bufferSize bytes
	paddedMessage := make([]byte, bufferSize)
	// Kopier meldingen inn i postkassen
	copy(paddedMessage, []byte(message))
	// Send meldingen til serveren
	_, err := conn.Write(paddedMessage)
	if err != nil {
		// Hvis det er en feil, si hva som gikk galt
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}

// Hovedfunksjonen for å håndtere alt klienten skal gjøre
func tcpFixedSizeClient(serverIP string, serverPort int) {
	const bufferSize = 1024 // Hvor stor hver melding skal være (her alltid 1024 bytes)

	// Prøv å koble til serveren
	conn, err := connectToServer(serverIP, serverPort)
	if err != nil {
		fmt.Println(err) // Skriv ut hva som gikk galt hvis det ikke fungerer
		return
	}
	defer conn.Close() // Når vi er ferdige, lukk forbindelsen automatisk

	// Motta velkomstmelding fra serveren
	welcomeMessage, err := receiveMessage(conn, bufferSize)
	if err != nil {
		fmt.Println(err) // Si hva som gikk galt
		return
	}
	fmt.Printf("Welcome message: %s\n", welcomeMessage) // Vis velkomstmeldingen

	// Les meldinger fra brukeren og send til serveren
	reader := bufio.NewReader(os.Stdin) // Lager en leser som kan lese hva du skriver i terminalen
	for {
		fmt.Print("Enter message (or 'exit' to quit): ") // Be brukeren skrive en melding
		message, _ := reader.ReadString('\n')            // Vent til brukeren trykker "Enter"
		message = message[:len(message)-1]               // Fjern linjeskiftet fra meldingen

		if message == "exit" { // Hvis brukeren skriver "exit", avslutt
			fmt.Println("Exiting...")
			return
		}

		// Send meldingen til serveren
		err := sendMessage(conn, message, bufferSize)
		if err != nil {
			fmt.Println(err) // Si hva som gikk galt
			continue         // Fortsett til neste runde
		}

		// Vent på svar fra serveren
		response, err := receiveMessage(conn, bufferSize)
		if err != nil {
			fmt.Println(err) // Si hva som gikk galt
			return
		}
		fmt.Printf("Server response: %s\n", response) // Skriv ut svaret fra serveren
	}
}

func main() {
	// Serverens IP-adresse og port
	serverIP := "10.100.23.204"
	serverPort := 34933

	tcpFixedSizeClient(serverIP, serverPort) // Start klienten
}
