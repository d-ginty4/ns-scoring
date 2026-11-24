export type Gender = "male" | "female";

type Common = {
	id: string
	createdAt: string
	updatedAt: string
	deletedAt: string
}

export type Region = Common & {
	name: string
}

export type Event = Common & {
	name: string
	eventNumber: number
	weightOptionsMen: number[]
	weightOptionsWomen: number[]
	scoringMethod: string
}

export type Division = Common & {
	name: string
	gender: Gender
}

export type Competition = Common & {
	name: string
	country: string
	iconName: string
	region: Region
}

export type Athlete = Common & {
	name: string
	country: string
	competition: Competition
	division: Division
}

export type Result = Common & {
	weight: number
	value: number
	didNotParticipate: boolean
	event: Event
	athlete: Athlete
}

export type Response<T> = {
	requestId: string
	data: T
	error: object
}
