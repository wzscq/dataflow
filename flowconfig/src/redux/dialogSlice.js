import { createSlice } from '@reduxjs/toolkit';

// Define the initial state using that type
const initialState = {
    dialogs:[]
}

export const dialogSlice = createSlice({
    name: 'dialog',
    initialState,
    reducers: {
        openDialog: (state,action) => {
            state.dialogs.push(action.payload);
        },
        closeDialog:(state,action) => {
            if(state.dialogs.length>0){
                delete state.dialogs[state.dialogs.length-1];
                state.dialogs=state.dialogs.filter(item=>item);
            }
        },
        updateDialogData:(state,action) => {
            if(state.dialogs.length>0){
                state.dialogs[state.dialogs.length-1].data=action.payload;
            }
        },
    }
});

// Action creators are generated for each case reducer function
export const { openDialog,closeDialog } = dialogSlice.actions

export default dialogSlice.reducer