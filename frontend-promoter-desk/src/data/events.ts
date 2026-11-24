import { request, resourcePath } from "./data"
import type { Event } from "../types"

const eventsPath = "/events"

export const getEvents = async (): Promise<Event[] | null> => {
	return request<Event[]>(eventsPath)
}

export const getEvent = async (id: number): Promise<Event | null> => {
	return request<Event>(resourcePath(eventsPath, id))
}
