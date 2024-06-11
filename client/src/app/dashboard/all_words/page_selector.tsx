import '../../../style/util.css'
import '../../../style/all_words.css'
import { v4 as uuidv4 } from 'uuid'

export interface PageSelectorInfo {
	page: number
	onPageChange: (newPage: number) => void
	pageCnt: number
}

export default function PageSelector(info: PageSelectorInfo) {
	const pageList = (() => {
		const list: Array<JSX.Element> = []
		if (info.pageCnt <= 5) {
			for (let i = 1; i <= info.pageCnt; i++) {
				list.push(
					<span
						key={uuidv4()}
						className='page-number'
						onClick={() => {
							info.onPageChange(i)
						}}
						style={
							info.page == i
								? {
										color: 'red',
										// eslint-disable-next-line no-mixed-spaces-and-tabs
								  }
								: {}
						}
					>
						{i}
					</span>,
				)
			}
		} else {
			for (let i = 1; i <= 3; i++) {
				list.push(
					<span
						key={uuidv4()}
						className='page-number'
						onClick={() => {
							info.onPageChange(i)
						}}
					>
						{i}
					</span>,
				)
			}
			list.push(<span key={uuidv4()}>...</span>)
			for (let i = info.pageCnt - 1; i <= info.pageCnt; i++) {
				list.push(
					<span
						key={uuidv4()}
						className='page-number'
						onClick={() => {
							info.onPageChange(i)
						}}
					>
						{i}
					</span>,
				)
			}
		}
		return list
	})()

	return (
		<div style={{ width: '100%', display: 'flex', justifyContent: 'center' }}>
			<button
				onClick={() => {
					if (info.page > 1) info.onPageChange(info.page - 1)
				}}
			>
				prev
			</button>
			<span
				style={{
					marginLeft: '0.2rem',
					marginRight: '0.2rem',
					display: 'flex',
					justifyContent: 'space-evenly',
				}}
				className='page-list-container'
			>
				{pageList}
			</span>
			<button
				onClick={() => {
					if (info.page < info.pageCnt) info.onPageChange(info.page + 1)
				}}
			>
				next
			</button>
		</div>
	)
}
