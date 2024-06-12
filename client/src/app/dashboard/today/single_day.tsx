import { Word } from './today'
import WordCard from './word_card'

export type Words = Array<Word>
export interface SingleDayInterface {
	day: number
	words: Words
}

export default function SingleDay(param: SingleDayInterface) {
	const words = param.words.map((value) => {
		return <WordCard key={value.id} word={value} />
	})

	return (
		<div id='single-day-root'>
			<div id='single-day-title'>
				<hr />
				<span style={{
                    marginLeft:'1rem',
                    marginRight:'1rem',
                }}>{param.day} days ago</span>
				<hr />
			</div>
			<div id='word-card-container'>{words}</div>
		</div>
	)
}
