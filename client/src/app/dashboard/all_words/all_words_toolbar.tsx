import '../../../style/util.css'
import '../../../style/all_words.css'
import { useState } from 'react'
import Detail, { DetailState } from '../detail'
import { generateURL, tryRequest } from '../../../util'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import config from '../../../config'

export type SortMethod = 'time' | 'alphabetic'

interface ToolBarInterface {
	method: SortMethod
	onMethodChange: (newMethod: SortMethod) => void
}

function ToolBar({ method, onMethodChange }: ToolBarInterface) {
	const navigate = useNavigate()
	const [detail, setDetail] = useState<DetailState>({
		show: false,
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
		submit: {
			text: 'add',
			onSubmit: async (word): Promise<string> => {
				if (word.word.length == 0) return "The length of the word can't be zero"
				for (const c of word.word) {
					if (!((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')))
						return 'The word contains invalid characters'
				}

				const result = { err: '' }
				await tryRequest(
					async () => {
						await axios.post(generateURL(config.api.word, '/word'), {
							method: 'add',
							parameters: word,
						})
						location.reload()
					},
					(err) => {
						result.err = 'something wrong'
						console.log(err)
					},
					navigate,
				)
				return result.err
			},
		},
	})
	const handleClose = () => setDetail({ ...detail, show: false })

	return (
		<div id='toolbar-root'>
			<button
				id='toolbar-add'
				onClick={() => {
					setDetail({ ...detail, show: true })
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
			{detail.show && (
				<Detail
					show={detail.show}
					word={detail.word}
					editable={detail.editable}
					submit={detail.submit}
					onClose={handleClose}
				/>
			)}
		</div>
	)
}

export default ToolBar
