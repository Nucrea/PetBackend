import classes from "./Content.module.css"

import { Tldraw } from 'tldraw'
import 'tldraw/tldraw.css'

export default function Content() {
    // return <div className={classes.content}>Test</div>
    // return <canvas className={classes.canvas}></canvas>

    return (
        <div className={classes.canvas}>
			<Tldraw />
		</div>
    );
}