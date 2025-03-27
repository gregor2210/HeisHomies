### PAKKE TAP BUGS
1) Fungere fint ved 25% pakketap. 50% faller egt alle heiser ut. MULIG vi kan redigere TimeOut tiden? Men vi må nokk rote med TCP instillinger
- Har gjort at den ikke kan prøve lenger en TimeOut er satt til. noe mer vi kan gjøre?

9) Ved høyt pakketap kan vi fort ende med at alle heisene er på men ikke klarer å komunisere med hverandre. Skal de ta ikke ta imot ordere? fordi det er ikke sikkert at det finnes en backup?


### Burde se på bugs
2) Når det er packetloss og en heis dc. Så kan en annen heis ta den orderen. så connecter de, synker ordre, så dc igjen og den andre heisen tar den orderen. blir en loop
13) Når en heis reconnecter så vil requstene dens lyse opp på panele. Men før den rekker å mota en ny WV så vil den sette lyset til det gamle
14) Dial og listen kan være i setup prosess når self går offline. Den vil da fortsatt connecte



### Ikke krise bugs
12) Priority value kalulering er fortsatt litt rart... Det fungerer nesten. men ikke helt.


### Formaliteter bugs
1) legg til Readme fil, gjør den bedre
2) Noen filer har ikke så mye komentering. Forekempel noe av orversatt heiskode. Legg til mer komentarer





--------------------------- Tror løst-----------------------------------------
