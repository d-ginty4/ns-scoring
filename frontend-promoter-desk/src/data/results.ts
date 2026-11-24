import { request, requestOk, resourcePath } from "./data"
import type { Result } from "../types"

export type ResultRequest = {
	weight: number
	value: number
	didNotParticipate: boolean
}

export type AthleteResultItem = ResultRequest & {
	eventId: number
}

export type AthleteResultsRequest = {
	results: AthleteResultItem[]
}

const resultsPath = (athleteId: number) => resourcePath("/results/athletes", athleteId)

const resultPath = (athleteId: number, eventId: number) => resourcePath(resultsPath(athleteId), "events", eventId)

export const getResults = async (athleteId: number): Promise<Result[] | null> => {
	return request<Result[]>(resultsPath(athleteId))
}

export const getResult = async (athleteId: number, eventId: number): Promise<Result | null> => {
	return request<Result>(resultPath(athleteId, eventId))
}

export const createResult = async (
	athleteId: number,
	eventId: number,
	result: ResultRequest,
): Promise<Result | null> => {
	return request<Result>(resultPath(athleteId, eventId), { method: "POST", body: result })
}

export const createAthleteResults = async (
	athleteId: number,
	results: AthleteResultsRequest,
): Promise<Result[] | null> => {
	return request<Result[]>(resultsPath(athleteId), { method: "POST", body: results })
}

export const updateResult = async (
	athleteId: number,
	eventId: number,
	result: ResultRequest,
): Promise<Result | null> => {
	return request<Result>(resultPath(athleteId, eventId), { method: "PUT", body: result })
}

export const deleteResult = async (athleteId: number, eventId: number): Promise<boolean> => {
	return requestOk(resultPath(athleteId, eventId), { method: "DELETE" })
}
