/**
 * @return {string}
 */
export function ProgressBar() {
    return (`
    <div class="progress">
        <div class="progress-bar bg-success" role="progressbar" 
            aria-valuenow="25" aria-valuemin="0" aria-valuemax="100"></div>
    </div>`)
}