import '../../../style/all_words.css'
import { useEffect, useLayoutEffect, useState } from 'react'
import { generateURL, tryRequest } from '../../../util'
import { NavigateFunction, useNavigate } from 'react-router-dom'
import axios from 'axios'
import config from '../../../config'
import ToolBar, { SortMethod } from './all_words_toolbar'
import WordLine from './word_line'
import PageSelector from './page_selector'

interface Word {
	id: string
	word: string
	explanation: string
	url: string
	created_at: Date
}

async function getWords(
	setWords: React.Dispatch<React.SetStateAction<Word[]>>,
	page: number,
	method: SortMethod,
	navigate: NavigateFunction,
	afterwards: () => void,
): Promise<string> {
	const errMessage = ''
	await tryRequest(
		async () => {
			const {
				data: { list },
			} = await axios.get<{ list: Array<Word> }>(
				generateURL(
					config.api.word,
					'/list/' + page.toString() + '?method=' + method,
				),
			)
			setWords(list)
			afterwards()
		},
		(error) => {
			console.log(error)
		},
		navigate,
	)
	return errMessage
}

function AllWords() {
	const [page, setPage] = useState(1)
	const [pageCnt, setPageCnt] = useState(0)
	const [sortMethod, setSortMethod] = useState<'time' | 'alphabetic'>(
		'alphabetic',
	)
	const [words, setWords] = useState<Array<Word>>([])
	const [errMessage, setErrMessage] = useState('')
	const navigate = useNavigate()

	useEffect(() => {
		tryRequest(
			async () => {
				const {
					data: { pages },
				} = await axios.get(generateURL(config.api.word, '/page'))
				setPageCnt(pages)
				console.log(pages)
			},
			(err) => {
				console.log(err)
			},
			navigate,
		)
	}, [navigate])
	useLayoutEffect(() => {
		(async () => {
			const err = await getWords(setWords, page, sortMethod, navigate, () => {})
			setErrMessage(err)
		})()
	}, [])
	// }, [page, sortMethod, navigate])

	const list = words.map((value, index) => {
		return (
			<div key={value.id}>
				<WordLine
					word={{
						id: value.id,
						index: config.pageLength * (page - 1) + index,
						word: value.word,
						explanation: value.explanation,
						url: value.url,
						created_at: value.created_at,
					}}
				/>
				<hr style={{ margin: '0', border: '0.3px solid grey' }} />
			</div>
		)
	})

	return (
		<div id='all_words_root'>
			<div id='toolbar'>
				<ToolBar
					method={sortMethod}
					onMethodChange={(newMethod: SortMethod) => {
						(async () => {
							const err = await getWords(
								setWords,
								page,
								newMethod,
								navigate,
								() => setSortMethod(newMethod),
							)
							setErrMessage(err)
						})()
					}}
				/>
				<hr />
			</div>
			<div>{errMessage}</div>
			<div id='main-content'>
				<div
					style={{
						display: 'flex',
						justifyContent: 'space-between',
						width: '100%',
						position: 'relative',
					}}
				>
					<span style={{ marginLeft: '1rem', width: '10%' }}>Number</span>
					<span style={{ width: '10%' }}>Word</span>
					<span
						style={{
							flexGrow: '5',
							maxWidth: '40%',
							textOverflow: 'ellipsis',
							overflow: 'hidden',
							marginRight: '2%',
						}}
					>
						Explanation
					</span>
					<span style={{ width: '10%' }}>Date</span>
					<span
						style={{
							marginLeft: '1rem',
							marginRight: '1rem',
							width: '15%',
							display: 'flex',
							justifyContent: 'center',
							alignItems: 'center',
						}}
					>
						Operations
					</span>
				</div>
				{list}
			</div>
			<div id='page-selector'>
				<PageSelector
					page={page}
					pageCnt={pageCnt}
					onPageChange={(newPage) => {
						(async () => {
							const err = await getWords(
								setWords,
								newPage,
								sortMethod,
								navigate,
								() => setPage(newPage),
							)
							setErrMessage(err)
						})()
					}}
				/>
			</div>
		</div>
	)
}

export default AllWords
