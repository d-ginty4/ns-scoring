import { request, resourcePath } from "./data"
import type { Division } from "../types"

const divisionsPath = "/divisions"

export const getDivisions = async (): Promise<Division[] | null> => {
	return request<Division[]>(divisionsPath)
}

export const getDivision = async (id: number): Promise<Division | null> => {
	return request<Division>(resourcePath(divisionsPath, id))
}
