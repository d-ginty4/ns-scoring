import { request, requestOk, resourcePath } from "./data"
import type { Competition } from "../types"

export type CompetitionRequest = {
	name: string
	country: string
	iconName?: string
	regionId: number
}

const competitionsPath = "/competitions"

export const getCompetitions = async (): Promise<Competition[] | null> => {
	return request<Competition[]>(competitionsPath)
}

export const getCompetition = async (id: number): Promise<Competition | null> => {
	return request<Competition>(resourcePath(competitionsPath, id))
}

export const createCompetition = async (competition: CompetitionRequest): Promise<Competition | null> => {
	return request<Competition>(competitionsPath, { method: "POST", body: competition })
}

export const updateCompetition = async (id: number, competition: CompetitionRequest): Promise<Competition | null> => {
	return request<Competition>(resourcePath(competitionsPath, id), { method: "PUT", body: competition })
}

export const deleteCompetition = async (id: number): Promise<boolean> => {
	return requestOk(resourcePath(competitionsPath, id), { method: "DELETE" })
}
