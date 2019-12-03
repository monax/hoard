package helpers

import (
	"github.com/monax/hoard/v7/meta"
)

type TestDocData struct {
	Type    string `json:"type"`
	RawData []byte `json:"data"`
}

type DocumentTest struct {
	Meta meta.Meta   `json:"meta"`
	Data TestDocData `json:"data"`
}

const LongText = `La testimonianza degli scrittori arabi al par che de’ diplomi cristiani
	/della Sicilia intorno Giorgio di Antiochia, conferma l’autorità civile
	delli ammiragli, che che si pensi de’ miei supposti su l’origine sua.
		Questa particolarità del diritto pubblico siciliano alla quale si è
	badato assai poco fin qui, ci aiuta a comprendere le vicissitudini
	dello Stato sotto i due Guglielmi, assai meglio che non faremmo col
	mero ordinamento dei sette grandi ufizii della Corona,[48] supponendo
	col Gregorio, che fosse stato fin da’ tempi di re Ruggiero qual si
	ritrae negli ultimi di Guglielmo il Buono, e che l’autorità di quegli
	ufizii si fosse estesa a tutti i sudditi, cristiani o musulmani. Erano
	gli elementi dell’azienda musulmana che tornavano a galla quando
	fu ristorata l’antica capitale. E dico delle istituzioni ed anco
	degli uomini. Guerrieri che avessero seguito in Terraferma il primo
	conte, uomini di mare, giuristi, segretarii, mercatanti, pedagoghi,
		camerieri; qual più qual meno caritatevoli, dissoluti e picchiapetto;
	bilingui e trilingui, barcheggianti tra due o tre religioni, versati
	nella letteratura arabica e nella scienza greca, dilettanti dell’arte
	bizantina e delle forme che prese in Siria, in Egitto o in Spagna: tali
	mi sembrano que’ Musulmani e Greci di Sicilia che la novella corte
	attirava, senza volerlo, nel castel di sopra di Palermo, insieme co’
	Levantini della tempra di Giorgio e coi prelati, i chierici e i nobili
	d’Italia e di Francia. Que’ costumi dissonanti s’armonizzaron pure un
	gran pezzo e produssero, nel corso del duodecimo secolo, due grandi
Statisti: orfani entrambi, maturati precocemente tra le agitazioni
	della corte di Palermo, somiglianti anco l’uno all’altro per tempra e
	cultura dell’intelletto, legislatori, buon massai, vaghi d’ogni scienza
	e filosofi più che cristiani: Ruggiero primo re e Federigo secondo
	imperatore; i due sultani battezzati di Sicilia, a’ quali l’Italia dee
	non piccola parte dell’incivilimento suo.`
