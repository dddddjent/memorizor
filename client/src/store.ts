import { configureStore } from '@reduxjs/toolkit'
import detailReducer from './app/dashboard/detail_slice'

export const store = configureStore({
    reducer: {
        detail: detailReducer,
    },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
