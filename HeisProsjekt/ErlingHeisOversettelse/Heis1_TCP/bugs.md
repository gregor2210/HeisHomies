### PAKKE TAP BUGS
1) Fungere fint ved 25% pakketap. 50% faller egt alle heiser ut. MULIG vi kan redigere TimeOut tiden? Men vi må nokk rote med TCP instillinger
- Har gjort at den ikke kan prøve lenger en TimeOut er satt til. noe mer vi kan gjøre?

9) Ved høyt pakketap kan vi fort ende med at alle heisene er på men ikke klarer å komunisere med hverandre. Skal de ta ikke ta imot ordere? fordi det er ikke sikkert at det finnes en backup?


### Burde se på bugs
11) Hvis heis er på vei en retning. Så stopper den på et mellomstopp. Da vil den ikke kunne ta en hall request i motsatt retning!


### Ikke krise bugs
12) Priority value kalulering er fortsatt litt rart... Det fungerer nesten. men ikke helt.

### Formaliteter bugs
1) legg til Readme fil, gjør den bedre
2) Noen filer har ikke så mye komentering. Forekempel noe av orversatt heiskode. Legg til mer komentarer



--------------------------- Tror løst-----------------------------------------
6) `LØST!!!` tror jeg . Lys skrur seg ikke på påtvers av heisene nå. Ønsker vi det? 

3) `TROR LØST`?? Programet går av og til inn i deadlock!!!!! Vet ikke hvor. Den siste meldingen som alltid skrives er faild: "Error sending, connection lost."

4) `TROR LØST`??? programmet havnet også en gang i deadlock etter "New state" ble skrevet idk

NR 3 OG 4 HAR JEG IKKE PRØVD Å LØSE ORDENTLIG!
Connected to localhost:8012
Setting ElevatorID: 1 to ONLINE!
2025/03/10 11:33:55 Elevator 0 is online
2025/03/10 11:33:55 Elevator 1 is online
2025/03/10 11:33:55 Elevator 2 is online
Returning from SetElevatorOnline
Starting handleReceive for elevator we are connected to
HANDLE RECEIVE STARTED, ID: 1
Error sending, connection lost.
^Csignal: interrupt

8) TROR LØST??? Sett heis til start poss må resette lampe

2) Ved mye packetloss så begynner TCP og bruke lenger tid på å retrye hvis den feiler flere ganger på rad. Ved høyt pakketap står den på DAIL veldig lenge.
- `Tror løst`. Den vil TimeOut etter hvor lang TimeOut er .

10) `LØST` legge til i lights at request sjekker om den er ulik forigje gang den sjekket. Slik at vi ikke sender unødvendig mange melinger til heis server?

5) `Løst`Vi kan ikke endre antall etasjer i heisen