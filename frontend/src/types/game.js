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
};
