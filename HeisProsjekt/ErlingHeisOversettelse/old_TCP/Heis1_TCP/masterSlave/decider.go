package masterSlave

type decisionType int

const (
	Primary decisionType = 1
	Backup  decisionType = -1
	Unknown decisionType = 0
)

var (
	decisionState = [3]decisionType{Unknown, Unknown, Unknown}
)

func SetDecisionState() {
	// Plan
	// If PC0 is online, it starts as the primary
	// If PC0 is offline, PC1 will be the primary.
	// If PC0 and PC1 are offline, PC2 will be the primary.

	// IF PC0 has been ofline and reconnects, it will be the new primary.
	// IF PC0 reconnects and reconnects faster to PC1 then PC2
	// -> Must check if PC1 sees PC2 as online,
	// ->-> if not, PC0 will be the primary, USIKKER PÅ OM VI VIL TA OVER ELLER IKKE ...
	// ->-> if yes, PC0 will be the primary

	// IF PC0 reconnects and reconnects faster to PC2 then PC1
	// -> Must chech if PC2 sees PC1 as online

	//IF PC1 has been O

}
