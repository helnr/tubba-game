export class EventMessage {
	constructor(type, payload) {
		this.type = type;
		this.payload = payload;
	}
}

export const EventMessageTypes = {
	ErrorEvent: "error_event",
	GameEvent: "game_event",
	JoinEvent: "join_event",
	LeaveEvent: "leave_event",
	ReadyEvent: "ready_event",
	ChangeTeamEvent: "change_team_event",
	StartEvent: "start_event",
	PlayedCardEvent: "played_card_event",
	TubbaEvent: "tubba_event",
};

export const STATUS = {
	LOADING: "loading",
	LOBBY: "lobby",
	READY: "ready",
	STARTED: "started",
	PAUSED: "paused",
	ENDED: "ended",
};
