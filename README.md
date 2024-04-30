# Diplomska-naloga-Vizualizacija-in-izvedba-verizne-replikacije
Diplomska naloga: Vizualizacija in izvedba verižne replikacije

Zagon programa:
go run main.go (default zagon s 5 strežniki)
go run main.go -n <št. strežnikov> (zagon s custom številom strežnikov)
OPOZORILO: zagon potrebuje nekaj sekund. Zaradi uporabe resta bo windows vprašal za dovoljenje izvajanja zaradi požarnega zidu 

REST klici:
(POST) localhost:8080/insert - vstavitev novega zapisa
JSON body example:
{
    "key": "string vrednost ključa",
    "value": "string vrednost zapisa"
}

(GET) localhost:8080/read - prebere vse zapise ali le izbranega
JSON body example za branje izbranega
{
    "command": "read",
    "arguments": ["server id", "string vrednost ključa"]
}

JSON body example za branje vseh
{
    "command": "readAll",
    "arguments": ["server id"]
}

(GET) localhost:8080/clock - en urin signal

št. urinih signalov potrebnih za dokončen zapis: 2 * število strežnikov
št. urinih signalov potrebnih za branje: 2

Podrobnejši zapisi za vstavljanje vrednosti se izpišejo v konzoli.
Nadzorno raven mislim implementirati v tekočem tednu.