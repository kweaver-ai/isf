import { createContext } from 'react'
import type { AppConfigContextType } from './type'

const AppConfigContext = createContext<AppConfigContextType>(null)

export default AppConfigContext