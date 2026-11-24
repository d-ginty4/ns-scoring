import { request, resourcePath } from "./data"
import type { Region } from "../types"

const regionsPath = "/regions"

export const getRegions = async (): Promise<Region[] | null> => {
	return request<Region[]>(regionsPath)
}

export const getRegion = async (id: number): Promise<Region | null> => {
	return request<Region>(resourcePath(regionsPath, id))
}
