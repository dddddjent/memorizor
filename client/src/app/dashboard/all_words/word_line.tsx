import { useNavigate } from 'react-router-dom'
import '../../../style/util.css'
import { useState } from 'react'
import Detail, { DetailState } from '../detail'
import { generateURL, tryRequest } from '../../../util'
import axios from 'axios'
import config from '../../../config'

export interface WordLineInfo {
	word: {
		index: number
		id: string
		word: string
		explanation: string
		url: string
		created_at: Date
	}
}

const MAX_EXPLANATION_LENGTH = 50

function WordLine(info: WordLineInfo) {
	const navigate = useNavigate()
	const updateDetail = {
		show: false,
		word: {
			word: info.word.word,
			explanation: info.word.explanation,
			url: info.word.url,
		},
		editable: {
			word: false,
			explanation: true,
			url: true,
		},
		submit: {
			text: 'update',
			onSubmit: async (word: {
				word: string
				explanation: string
				url: string
			}): Promise<string> => {
				const result = { err: '' }
				await tryRequest(
					async () => {
						await axios.post(generateURL(config.api.word, '/word'), {
							method: 'update',
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
	}
	const [detail, setDetail] = useState<DetailState>(updateDetail)
	const handleClose = () => setDetail({ ...detail, show: false })
	const previewDetail = {
		show: false,
		word: {
			word: info.word.word,
			explanation: info.word.explanation,
			url: info.word.url,
		},
		editable: {
			word: false,
			explanation: false,
			url: false,
		},
		submit: {
			text: 'OK',
			onSubmit: async (): Promise<string> => {
				handleClose()
				return ''
			},
		},
	} as DetailState

	const handleDelete = () => {
		tryRequest(
			async () => {
				await axios.post(generateURL(config.api.word, '/word'), {
					method: 'delete',
					parameters: { id: info.word.id },
				})
				location.reload()
			},
			(err) => {
				console.log(err)
			},
			navigate,
		)
	}

	return (
		<div id='word-line-container'>
			<span
				style={{ marginLeft: '1rem', width: '10%' }}
				className='word-line-span'
				onClick={() => setDetail({ ...previewDetail, show: true })}
			>
				{info.word.index + 1}
			</span>
			<span
				style={{ width: '10%' }}
				className='word-line-span'
				onClick={() => setDetail({ ...previewDetail, show: true })}
			>
				{info.word.word}
			</span>
			<span
				style={{
					flexGrow: '5',
					maxWidth: '40%',
					textOverflow: 'ellipsis',
					overflow: 'hidden',
					marginRight: '2%',
				}}
				className='word-line-span'
				onClick={() => setDetail({ ...previewDetail, show: true })}
			>
				{info.word.explanation.length > MAX_EXPLANATION_LENGTH
					? info.word.explanation.slice(0, MAX_EXPLANATION_LENGTH) + '...'
					: info.word.explanation}
			</span>
			<span
				style={{ width: '10%' }}
				className='word-line-span'
				onClick={() => setDetail({ ...previewDetail, show: true })}
			>
				{new Date(info.word.created_at).toLocaleDateString('en-US')}
			</span>
			<span
				style={{
					marginLeft: '1rem',
					marginRight: '1rem',
					width: '15%',
					display: 'flex',
					justifyContent: 'center',
					alignItems: 'center',
				}}
				className='word-line-span'
			>
				<button
					onClick={() => {
						setDetail({ ...updateDetail, show: true })
					}}
				>
					update
				</button>
				<button style={{ marginLeft: '0.5rem' }} onClick={() => handleDelete()}>
					delete
				</button>
				{detail.show && (
					<Detail
						show={detail.show}
						word={detail.word}
						editable={detail.editable}
						submit={detail.submit}
						onClose={handleClose}
					/>
				)}
			</span>
		</div>
	)
}

export default WordLine
