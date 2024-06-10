import '../../../style/util.css'
import '../../../style/all_words.css'
import { useAppDispatch } from '../../../util'
import { open as openDetail } from '../detail_slice'

export type SortMethod = 'time' | 'alphabetic'

interface ToolBarInterface {
    method: SortMethod
    onMethodChange: (newMethod: SortMethod) => void
}

function ToolBar({ method, onMethodChange }: ToolBarInterface) {
    const dispatch = useAppDispatch()

    return (
        <div id='toolbar-root'>
            <button
                id='toolbar-add'
                onClick={() => {
                    dispatch(
                        openDetail({
                            show: true,
                            word: {
                                word: '',
                                explanation: '',
                                url: '',
                            },
                            editable: {
                                word: true,
                                explanation: true,
                                url: true,
                            },
                        }),
                    )
                }}
            >
                Add
            </button>
            <div id='toolbar-search'>
                <label id='toolbar-search-label'>Search a word: </label>
                <input id='toolbar-search-input' />
            </div>
            <select
                id='toolbar-select'
                value={method}
                onChange={(e) => {
                    onMethodChange(e.target.value as SortMethod)
                }}
            >
                <option value='time'>Time</option>
                <option value='alphabetic'>Alphabetic</option>
            </select>
        </div>
    )
}

export default ToolBar
