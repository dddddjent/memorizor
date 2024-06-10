import { createSlice, PayloadAction } from '@reduxjs/toolkit'

import { DetailInterface } from './detail'
import { RootState } from '../../store'

export type DetailState = Omit<DetailInterface, 'onClose'>

const initialState: DetailState = {
    show: false,
    word: {
        word: 'A',
        explanation: 'None',
        url: 'www.example.com',
    },
    editable: {
        word: false,
        explanation: false,
        url: false,
    },
}

export const detailSlice = createSlice({
    name: 'detail',
    initialState,
    reducers: {
        open: (state, action: PayloadAction<DetailState>) => {
            state.show = true
            state.word = action.payload.word
            state.editable = action.payload.editable
        },
        close: (state) => {
            state.show = false
        },
    },
})

export const { open, close } = detailSlice.actions
export const selectDetail = (state: RootState) => state.detail
export default detailSlice.reducer
