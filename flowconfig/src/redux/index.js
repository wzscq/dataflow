import { configureStore } from '@reduxjs/toolkit'

import frameReducer from './frameSlice';
import flowReducer from './flowSlice';
import i18nReducer from './i18nSlice';
import dialogReducer from './dialogSlice';
import debugReducer from './debugSlice';

export default configureStore({
  reducer: {
    frame:frameReducer,
    flow:flowReducer,
    i18n:i18nReducer,
    dialog:dialogReducer,
    debug:debugReducer
  }
});