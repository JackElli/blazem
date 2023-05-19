import { c as create_ssr_component, b as add_attribute } from "../../../../chunks/index2.js";
const _page_svelte_svelte_type_style_lang = "";
const css = {
  code: "#change_log.svelte-16qhpbc h1.svelte-16qhpbc{font-size:x-large}#change_log.svelte-16qhpbc h1.svelte-16qhpbc{font-size:x-large}ul.svelte-16qhpbc.svelte-16qhpbc{list-style-type:circle}li.svelte-16qhpbc span.svelte-16qhpbc{font-weight:500}#help.svelte-16qhpbc h1.svelte-16qhpbc{font-size:x-large;font-weight:500}",
  map: null
};
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  $$result.css.add(css);
  return `${$$result.head += `<!-- HEAD_svelte-1gr2w2r_START -->${$$result.title = `<title>Blazem | Help</title>`, ""}<!-- HEAD_svelte-1gr2w2r_END -->`, ""}
<div class="flex gap-2 mt-4">
    <h1${add_attribute("class", `${"text-[#3d3d75] underline"} hover:underline cursor-pointer`, 0)}>Help
    </h1>
    
    <h1${add_attribute(
    "class",
    `${""} hover:underline cursor-pointer`,
    0
  )}>Changelog
    </h1>
    
    <h1${add_attribute(
    "class",
    `${""} hover:underline cursor-pointer`,
    0
  )}>Feature Requests
    </h1></div>
<div class="mt-4">${`<div id="help" class="w-1/2 svelte-16qhpbc"><h1 class="svelte-16qhpbc">What is Blazem?</h1>
            <p>Blazem is a <bold>Data Retention Service</bold>, created to help
                learn GO and the other aspects of software development. It has
                taught the fundementals of <bold>Distributed Systems</bold> and
                has been heavily inspired by <bold>Couchbase</bold>.
            </p>
            <h1 class="svelte-16qhpbc">How to use Blazem.</h1>
            <p>When you first log in to Blazem, you&#39;ll find yourself on the
                folders page. This gives a quick overview of which folders you
                have within Blazem.
            </p>
            <ul class="svelte-16qhpbc"><li>Click on the folder <bold>text</bold> to go into the folder.
                </li></ul>

            <p>Once you are in the folder, you can see all of the documents
                currently stored within that folder. The <bold>document key</bold>
                is located in the header and the <bold>document fields</bold> are
                in the body.
            </p>
            <ul class="svelte-16qhpbc"><li>Add some data by clicking the <bold>Add Data</bold> button in
                    the top right.
                </li>
                <li>This will bring up the Add Data modal.</li>
                <ol><li>Choose a folder for the data to reside (the current
                        folder is selected by default)
                    </li>
                    <li>Choose a key, these are auto-generated.</li>
                    <li>Type in your value or upload an image.</li></ol></ul>

            <p>You can find your way around using the Header bar to access:</p>
            <ul class="svelte-16qhpbc"><li>Folders - this brings you back to your folder view.</li>
                <li>Query - able to search and query your data.</li>
                <li>Nodes - a quick overview of the health of your nodes.</li>
                <li>Rules - a set of tasks created to run at specific times.
                </li></ul>
            <h1 class="svelte-16qhpbc">Is Blazem free to use?</h1>
            <p>Blazem is totally <bold>free to use</bold>, with the following
                limits:
            </p>
            <ul class="svelte-16qhpbc"><li>Up to 2 nodes, is totally free.</li>
                <li>Up to 3 rules deployed, is totally free.</li></ul>
            <p>To increase these limits, <bold>Blazem Premium</bold> will be required.
            </p>

            <h1 class="svelte-16qhpbc">How do I query my data?</h1>
            <p>Querying your data is a quick way of adding some logic to your
                searches. For instance, you want to find all the videos that are
                about the beach. Easy. Define a query that matches those
                criterea.
            </p>
            <ul class="svelte-16qhpbc"><li>Navigate to the <bold>Query Tab</bold> located in the header.
                </li>
                <li>Type a query in the query box: <bold>SELECT all</bold></li>
                <li>Hit <bold>Execute</bold> and watch your documents appear.
                </li></ul>
            <h1 class="svelte-16qhpbc">What is a Rule?</h1>
            <p>A rule is a set of tasks to be performed at certain times of the
                day specified by you. For example you could create a Rule to
                SELECT all documents WHERE key is LIKE &quot;todo&quot; and export them to
                Couchbase at 12 AM every day.
            </p>
            <h1 class="svelte-16qhpbc">What if one of my nodes goes down?</h1>
            <ul class="svelte-16qhpbc"><li>If one of the follower nodes goes down, it doesn&#39;t matter,
                    computers sometimes fail. Just reconnect when you can to
                    take full advantage of <bold>High Availability</bold></li>
                <li>If the master node goes down, all the data will be pushed to
                    the <bold>next in line</bold> (as long as you have more than
                    one node connected). No data will be lost.
                </li>
                <li>If a node goes down, it will reconnect automatically, unless
                    you remove it by clicking on the node when <bold>red</bold> or
                    10 minutes has elapsed and the node is still down.
                </li></ul>

            <hr>
            <h1 class="svelte-16qhpbc">Why Blazem?</h1>
            <p>Blazem is a memory first, <bold>Distributed System Architecture</bold>, allowing the <bold>rapid</bold> retrieval of documents and
                peace of mind, knowing your documents will be safe. Safety is
                our main priority, saving data to disk, to memory and across
                multiple nodes, ensuring
                <bold>High Availability</bold> and data safety when you need it most.
            </p>

            <h1 class="svelte-16qhpbc">How to deploy?</h1>
            <p>Blazem doesn&#39;t require any <bold>specialised</bold> hardware. Instead,
                run it on any number of Raspberry PIs/ any other hardware you like.
            </p>
            <ol><li>Install Blazem on your device</li>
                <li>Run Blazem</li>
                <li>Open Web UI at IP:3100</li>
                <li>Upload documents and away you go</li></ol>

            <h1 class="svelte-16qhpbc">What type of documents does Blazem support?</h1>
            <p>Blazem stores documents in the <bold>JSON</bold> format.</p>
            <p>Supports:</p>
            <ul class="svelte-16qhpbc"><li>Simple text</li>
                <li>Images</li>
                <li>Files (able to peek and download)</li>
                <li>Video (stored on disk)</li>
                <li>Multi-field JSON</li></ul></div>`}
</div>`;
});
export {
  Page as default
};
