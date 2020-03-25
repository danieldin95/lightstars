/**
 * @return {string}
 */
export function ProgressBar(id) {
    return (`
    <div id="${id}" class="progress">
        <div class="progress-bar bg-success" role="progressbar" 
            aria-valuenow="25" aria-valuemin="0" aria-valuemax="100">
        </div>
    </div>`)
}