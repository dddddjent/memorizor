import { useState } from 'react'
import { leftAlignedPosition } from '../../util'

export interface DetailInterface {
	show: boolean
	word: {
		word: string
		explanation: string
		url: string
	}
	editable: {
		word: boolean
		explanation: boolean
		url: boolean
	}
	submit: {
		text: string
		onSubmit: (word: {
			word: string
			explanation: string
			url: string
		}) => Promise<string>
	}
	onClose: () => void
}
export type DetailState = Omit<DetailInterface, 'onClose'>

const LABEL_LEFT = 4

export default function Detail(param: DetailInterface) {
	const [content, setContent] = useState({
		word: param.word.word,
		explanation: param.word.explanation,
		url: param.word.url,
	})
	const [errMsg, setErrMsg] = useState('')

	const handleSubmit = async () => {
		const result = await param.submit.onSubmit({
			word: content.word,
			explanation: content.explanation,
			url: content.url,
		})
		setErrMsg(result)
		if (result == '') param.onClose()
	}

	return (
		<div id='detail-root'>
			<div id='detail-content-container'>
				<div style={{ width: '100%', display: 'flex' }}>
					<button id='detail-close' onClick={param.onClose}>
						close
					</button>
				</div>
				<div id='detail-word'>
					{param.editable.word && (
						<input
							style={{
								...leftAlignedPosition(LABEL_LEFT),
								width: '70%',
								overflow: 'auto',
								fontSize: '2.3rem',
							}}
							type='text'
							placeholder='Your new word'
							value={content.word}
							onInput={(e: React.FormEvent<HTMLInputElement>) => {
								setContent({
									...content,
									word: (e.target as HTMLInputElement).value,
								})
							}}
						/>
					)}
					{!param.editable.word && (
						<span
							style={{
								...leftAlignedPosition(LABEL_LEFT),
								width: '70%',
								overflow: 'auto',
							}}
						>
							{param.word.word}
						</span>
					)}
				</div>
				<div id='detail-content'>
					<div id='detail-explanation-row'>
						<span
							style={{
								marginLeft: '2rem',
								width: '6rem',
							}}
						>
							Explanation:
						</span>
						{param.editable.explanation && (
							<textarea
								style={{
									overflowY: 'auto',
									overflowWrap: 'anywhere',
									background: '#dddddd',
									flexGrow: '100',
									marginLeft: '1rem',
									fontSize: '1rem',
								}}
								placeholder="Word's explanation"
								value={content.explanation}
								onInput={(e: React.FormEvent<HTMLTextAreaElement>) => {
									setContent({
										...content,
										explanation: (e.target as HTMLTextAreaElement).value,
									})
								}}
							/>
						)}
						{!param.editable.explanation && (
							<span
								style={{
									overflowY: 'auto',
									overflowWrap: 'anywhere',
									background: '#dddddd',
									flexGrow: '100',
									marginLeft: '1rem',
								}}
							>
								{param.word.explanation}
							</span>
						)}
					</div>
					<div id='detail-url-row'>
						<span
							style={{
								marginLeft: '2rem',
								width: '6rem',
								minWidth: '6rem',
							}}
						>
							URL:
						</span>
						{param.editable.url && (
							<input
								style={{
									overflowY: 'auto',
									textOverflow: 'ellipsis',
									flexGrow: '100',
									marginLeft: '1rem',
									marginTop: '0.2rem',
									marginBottom: '0.2rem',
								}}
								type='text'
								placeholder="Word's external url"
								value={content.url}
								onInput={(e: React.FormEvent<HTMLInputElement>) => {
									setContent({
										...content,
										url: (e.target as HTMLInputElement).value,
									})
								}}
							/>
						)}
						{!param.editable.url && (
							<span
								style={{
									overflowY: 'auto',
									textOverflow: 'ellipsis',
									flexGrow: '100',
									marginLeft: '1rem',
								}}
							>
								{param.word.url}
							</span>
						)}
					</div>
					<div id='detail-bottom-row'>
						<span
							style={{
								...leftAlignedPosition(50),
								transform: 'translate(-50%, -50%)',
								color: 'red',
								fontSize: '0.7rem',
							}}
						>
							{errMsg}
						</span>
						<button
							style={{
								marginLeft: 'auto',
							}}
							onClick={handleSubmit}
						>
							{param.submit.text}
						</button>
					</div>
				</div>
			</div>
		</div>
	)
}
