import '../../../style/util.css'
import { leftAlignedPosition } from '../../../util.ts'

export interface WordLineInfo {
    word: {
        index: number
        id: string
        word: string
        explanation: string
        created_at: Date
    }
}

const MAX_EXPLANATION_LENGTH = 10

function WordLine(info: WordLineInfo) {
    return (
        <div id='word-line-container'>
            <span style={{ marginLeft: '1rem' }}>{info.word.index + 1}</span>
            <span style={leftAlignedPosition(10)}>{info.word.word}</span>
            <span style={leftAlignedPosition(30)}>
                {info.word.explanation.length > MAX_EXPLANATION_LENGTH
                    ? info.word.explanation.slice(0, MAX_EXPLANATION_LENGTH) + '...'
                    : info.word.explanation}
            </span>
            <span style={{ marginRight: '1rem' }}>
                {new Date(info.word.created_at).toLocaleDateString('en-US')}
            </span>
        </div>
    )
}

export default WordLine
