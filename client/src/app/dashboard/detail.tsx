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
				<div style={{ width: '100%', display: 'flex', marginTop: '0.2rem' }}>
					<button
						id='detail-close'
						onClick={param.onClose}
						style={{ width: '3rem' }}
					>
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
								background: 'inherit',
								borderRadius: '6px',
								borderStyle: 'solid',
								borderColor: '#eeeeee',
								borderWidth: '1px',
								padding: '0.5rem',
								color: 'inherit',
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
									background: 'inherit',
									borderRadius: '6px',
									borderStyle: 'solid',
									borderColor: '#eeeeee',
									borderWidth: '1px',
									color: 'inherit',
									padding: '0.5rem',
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
									flexGrow: '100',
									marginLeft: '1rem',

									borderRadius: '10px',
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
									marginTop: '0.4rem',
									marginBottom: '0.4rem',

									background: 'inherit',
									borderRadius: '6px',
									borderStyle: 'solid',
									borderColor: '#eeeeee',
									borderWidth: '1px',
									paddingLeft: '0.5rem',
									color: 'inherit',
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
									marginTop: '0.4rem',
									marginBottom: '0.4rem',
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
								width: '3rem',
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
