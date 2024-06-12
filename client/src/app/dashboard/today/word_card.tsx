import '../../../style/util.css'
import '../../../style/today.css'
import { Word } from './today'
import Detail, { DetailState } from '../detail'
import { useState } from 'react'
import { generateURL, tryRequest } from '../../../util'
import axios from 'axios'
import config from '../../../config'
import { useNavigate } from 'react-router-dom'

export interface WordCardInterface {
	word: Word
}

export default function WordCard(param: WordCardInterface) {
	const navigate = useNavigate()
	const [detail, setDetail] = useState<DetailState>({
		show: false,
		word: {
			word: param.word.word,
			explanation: param.word.explanation,
			url: param.word.url,
		},
		editable: {
			word: false,
			explanation: false,
			url: false,
		},
		submit: {
			text: 'OK',
			onSubmit: async (): Promise<string> => {
				return ''
			},
		},
	})

	const clickedToday = (() => {
		const today = new Date()
		const clicked_at = new Date(param.word.clicked_at)
		if (today.getFullYear() !== clicked_at.getFullYear()) return false
		if (today.getMonth() !== clicked_at.getMonth()) return false
		if (today.getDay() !== clicked_at.getDay()) return false
		return true
	})()

	const handleClose = () => setDetail({ ...detail, show: false })
	const handleDetail = () => {
		setDetail({
			...detail,
			show: true,
			submit: {
				text: 'OK',
				onSubmit: async (): Promise<string> => {
					if (clickedToday === true) {
                        handleClose()
						return ''
					}
					const result = { err: '' }
					await tryRequest(
						async () => {
							await axios.post(generateURL(config.api.word, '/word'), {
								method: 'click',
								parameters: {
									id: param.word.id,
								},
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
	}

	return (
		<div
			id='word-card-root'
			style={{
				backgroundColor: clickedToday ? '#ddffdd' : 'grey',
			}}
		>
			<div
				style={{
					marginLeft: '0.5rem',
					marginRight: '0.5rem',
					display: 'flex',
					alignItems: 'center',
					justifyContent: 'center',
					flexGrow: '1',
				}}
			>
				<p
					style={{
						fontSize: '1.2rem',
						textOverflow: 'ellipsis',
						overflow: 'hidden',
						maxWidth: '100%',
					}}
				>
					{param.word.word}
				</p>
			</div>
			<div
				style={{
					width: '100%',
					display: 'flex',
					marginBottom: '0.2rem',
				}}
			>
				<button
					style={{
						marginLeft: 'auto',
						marginRight: '0.2rem',
					}}
					onClick={handleDetail}
				>
					detail
				</button>
			</div>
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
