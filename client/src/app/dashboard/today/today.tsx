import { useEffect, useState } from 'react'
import '../../../style/today.css'
import '../../../style/util.css'
import Toolbar from './toolbar'
import { generateURL, tryRequest } from '../../../util'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import config from '../../../config'
import SingleDay from './single_day'
import { v4 } from 'uuid'

export interface Word {
	id: string
	word: string
	explanation: string
	url: string
	clicked_at: string | Date
}

export type TodayList = Array<Array<Word>>

export default function Today() {
	const [todayList, setTodayList] = useState<TodayList>([[]])
	const [refreshing, setRefreshing] = useState(false)
	const navigate = useNavigate()
	const refresh = () => {
		setRefreshing(true)
		tryRequest(
			async () => {
				// TS bug: today's word clicked_at is not Date! It's a string
				const {
					data: { today },
				}: { data: { today: TodayList } } = await axios.get(
					generateURL(config.api.word, '/today'),
				)
				setTodayList(today)
				setRefreshing(false)
			},
			() => {},
			navigate,
		)
	}

	useEffect(() => {
		refresh()
	}, [])

	const fib = [0, 1]
	const days = todayList.map((value: Array<Word>, index) => {
		fib[index % 2] = fib[index % 2] + fib[(index + 1) % 2]
		return <SingleDay key={v4()} day={fib[index % 2]} words={value}></SingleDay>
	})

	return (
		<div id='today_root'>
			<div id='toolbar'>
				<Toolbar refreshing={refreshing} onRefresh={() => refresh()} />
			</div>
			<div id='today-main-content'>{days}</div>
		</div>
	)
}
