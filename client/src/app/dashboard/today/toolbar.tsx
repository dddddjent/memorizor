export interface ToolbarInterface {
	refreshing: boolean
	onRefresh: () => void
}

export default function Toolbar(param: ToolbarInterface) {
	return (
		<div id='toolbar-root'>
			<div id='toolbar-title'>Words Today</div>
			<button
				style={{ marginRight: '1rem', marginLeft: 'auto' }}
				disabled={param.refreshing}
				onClick={param.onRefresh}
			>
				Refresh
			</button>
		</div>
	)
}
