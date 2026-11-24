import { request, requestOk, resourcePath } from "./data"
import type { Athlete } from "../types"

export type AthleteRequest = {
	name: string
	country: string
	divisionId: number
}

export type AthletesRequest = {
	athletes: AthleteRequest[]
}

const athletesPath = (competitionId: number) => resourcePath("/competitions", competitionId, "athletes")

export const getAthletes = async (competitionId: number): Promise<Athlete[] | null> => {
	return request<Athlete[]>(athletesPath(competitionId))
}

export const getAthlete = async (competitionId: number, athleteId: number): Promise<Athlete | null> => {
	return request<Athlete>(resourcePath(athletesPath(competitionId), athleteId))
}

export const createAthlete = async (competitionId: number, athlete: AthleteRequest): Promise<Athlete | null> => {
	return request<Athlete>(athletesPath(competitionId), { method: "POST", body: athlete })
}

export const createAthletes = async (competitionId: number, athletes: AthletesRequest): Promise<Athlete[] | null> => {
	return request<Athlete[]>(`${athletesPath(competitionId)}/bulk`, { method: "POST", body: athletes })
}

export const updateAthlete = async (
	competitionId: number,
	athleteId: number,
	athlete: AthleteRequest,
): Promise<Athlete | null> => {
	return request<Athlete>(resourcePath(athletesPath(competitionId), athleteId), { method: "PUT", body: athlete })
}

export const deleteAthlete = async (competitionId: number, athleteId: number): Promise<boolean> => {
	return requestOk(resourcePath(athletesPath(competitionId), athleteId), { method: "DELETE" })
}
