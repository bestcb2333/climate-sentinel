export interface User {
  id: number
  createdAt: Date
  updatedAt: Date
  name: string
  email: string
  admin: boolean
  region: Region | null
  notices: Notice[]
  events: Event[]
  messages: Message[]
}

export interface Region {
  id: number
  updatedAt: Date
  name: string
  description: string
  coordinate: [number, number][][][]
  users: User[]
  events: Event[]
  histories: History[]
  resources: Resource[]
  routes: Route[]
}

export interface Event {
  id: number
  createdAt: Date
  name: string
  type: string
  regions: Region[]
  startTime: Date
  endTime: Date
  user: User | null
  severity: number
  coordinate: [number, number]
  description: string
}

export interface History {
  id: number
  createdAt: Date
  region: Region
  type: string
  time: Date
  maxTemperature: number
  minTemperature: number
  axgTemperature: number
  windSpeed: number
  visibility: number
  rainfall: number | null
  severity: number
  source: string
}

export interface Resource {
  id: number
  updatedAt: Date
  type: string
  name: string
  quantity: number
  region: Region
  coordinate: [number, number]
  available: boolean
}

export interface Wiki {
  id: number
  name: string
  notices: Notice[]
}

export interface Notice {
  id: number
  createdAt: Date
  updatedAt: Date
  user: User | null
  title: string
  content: string
}

export interface Route {
  id: number
  createdAt: Date
  type: string,
  name: string
  coordinate: [number, number][][]
  description: string
  available: boolean
  rate: number
  region: Region
}

export interface Message {
  id: number
  createdAt: Date
  user: User
  role: string
  content: string
}
