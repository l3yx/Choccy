export function convertHTML(str) {
    const characters = [/&/g, /</g, />/g, /"/g, /'/g,/\n/g,/ /g];
    const entities = ["&amp;", "&lt;", "&gt;", "&quot;", "&apos;","<br>","&nbsp;"];
    for (let i = 0; i < characters.length; i++) {
        str = str.replace(characters[i], entities[i]);
    }
    return str;
}