<script>
    import PageTitle from "$lib/components/PageTitle.svelte";

    let tabSelected = "help";
</script>

<svelte:head>
    <title>Blazem | Help</title>
</svelte:head>
<div class="flex items-center">
    <PageTitle>Help</PageTitle>
</div>
<div class="flex gap-2 mt-4">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <h1
        class={`${
            tabSelected == "help" ? "text-[#3d3d75] underline" : ""
        } hover:underline cursor-pointer`}
        on:click={() => (tabSelected = "help")}
    >
        Help
    </h1>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <h1
        class={`${
            tabSelected == "changelog" ? "text-[#3d3d75] underline" : ""
        } hover:underline cursor-pointer`}
        on:click={() => (tabSelected = "changelog")}
    >
        Changelog
    </h1>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <h1
        class={`${
            tabSelected == "features" ? "text-[#3d3d75] underline" : ""
        } hover:underline cursor-pointer`}
        on:click={() => (tabSelected = "features")}
    >
        Feature Requests
    </h1>
</div>
<div class="mt-4">
    {#if tabSelected == "help"}
        <div id="help" class="w-1/2">
            <h1>What is Blazem?</h1>
            <p>
                Blazem is a <bold>Data Retention Service</bold>, created to help
                learn GO and the other aspects of software development. It has
                taught the fundementals of <bold>Distributed Systems</bold> and
                has been heavily inspired by <bold>Couchbase</bold>.
            </p>
            <h1>How to use Blazem.</h1>
            <p>
                When you first log in to Blazem, you'll find yourself on the
                folders page. This gives a quick overview of which folders you
                have within Blazem.
            </p>
            <ul>
                <li>
                    Click on the folder <bold>text</bold> to go into the folder.
                </li>
            </ul>

            <p>
                Once you are in the folder, you can see all of the documents
                currently stored within that folder. The <bold
                    >document key</bold
                >
                is located in the header and the <bold>document fields</bold> are
                in the body.
            </p>
            <ul>
                <li>
                    Add some data by clicking the <bold>Add Data</bold> button in
                    the top right.
                </li>
                <li>This will bring up the Add Data modal.</li>
                <ol>
                    <li>
                        Choose a folder for the data to reside (the current
                        folder is selected by default)
                    </li>
                    <li>Choose a key, these are auto-generated.</li>
                    <li>Type in your value or upload an image.</li>
                </ol>
            </ul>

            <p>You can find your way around using the Header bar to access:</p>
            <ul>
                <li>Folders - this brings you back to your folder view.</li>
                <li>Query - able to search and query your data.</li>
                <li>Nodes - a quick overview of the health of your nodes.</li>
                <li>
                    Rules - a set of tasks created to run at specific times.
                </li>
            </ul>
            <h1>Is Blazem free to use?</h1>
            <p>
                Blazem is totally <bold>free to use</bold>, with the following
                limits:
            </p>
            <ul>
                <li>Up to 2 nodes, is totally free.</li>
                <li>Up to 3 rules deployed, is totally free.</li>
            </ul>
            <p>
                To increase these limits, <bold>Blazem Premium</bold> will be required.
            </p>

            <h1>How do I query my data?</h1>
            <p>
                Querying your data is a quick way of adding some logic to your
                searches. For instance, you want to find all the videos that are
                about the beach. Easy. Define a query that matches those
                criterea.
            </p>
            <ul>
                <li>
                    Navigate to the <bold>Query Tab</bold> located in the header.
                </li>
                <li>Type a query in the query box: <bold>SELECT all</bold></li>
                <li>
                    Hit <bold>Execute</bold> and watch your documents appear.
                </li>
            </ul>
            <h1>What is a Rule?</h1>
            <p>
                A rule is a set of tasks to be performed at certain times of the
                day specified by you. For example you could create a Rule to
                SELECT all documents WHERE key is LIKE "todo" and export them to
                Couchbase at 12 AM every day.
            </p>
            <h1>What if one of my nodes goes down?</h1>
            <ul>
                <li>
                    If one of the follower nodes goes down, it doesn't matter,
                    computers sometimes fail. Just reconnect when you can to
                    take full advantage of <bold>High Availability</bold>
                </li>
                <li>
                    If the master node goes down, all the data will be pushed to
                    the <bold>next in line</bold> (as long as you have more than
                    one node connected). No data will be lost.
                </li>
                <li>
                    If a node goes down, it will reconnect automatically, unless
                    you remove it by clicking on the node when <bold>red</bold> or
                    10 minutes has elapsed and the node is still down.
                </li>
            </ul>

            <hr />
            <h1>Why Blazem?</h1>
            <p>
                Blazem is a memory first, <bold
                    >Distributed System Architecture</bold
                >, allowing the <bold>rapid</bold> retrieval of documents and
                peace of mind, knowing your documents will be safe. Safety is
                our main priority, saving data to disk, to memory and across
                multiple nodes, ensuring
                <bold>High Availability</bold> and data safety when you need it most.
            </p>

            <h1>How to deploy?</h1>
            <p>
                Blazem doesn't require any <bold>specialised</bold> hardware. Instead,
                run it on any number of Raspberry PIs/ any other hardware you like.
            </p>
            <ol>
                <li>Install Blazem on your device</li>
                <li>Run Blazem</li>
                <li>Open Web UI at IP:3100</li>
                <li>Upload documents and away you go</li>
            </ol>

            <h1>What type of documents does Blazem support?</h1>
            <p>Blazem stores documents in the <bold>JSON</bold> format.</p>
            <p>Supports:</p>
            <ul>
                <li>Simple text</li>
                <li>Images</li>
                <li>Files (able to peek and download)</li>
                <li>Video (stored on disk)</li>
                <li>Multi-field JSON</li>
            </ul>
        </div>
    {:else if tabSelected == "changelog"}
        <div id="change_log">
            <h1>0.1.12</h1>
            <ul>
                <li>Added more backend tests</li>
                <li>
                    Updated button design and updated breadcrumb design to fit
                    new theme
                </li>
                <li>Changed design of modals</li>
                <li>
                    Thinking about adding a 'Check for updates' section so we
                    can update more easily
                </li>
            </ul>
            <h1>0.1.11</h1>
            <ul>
                <li>
                    Backend changes again, now make use of a Register function
                    to register all of the endpoints
                </li>
                <li>This will make room for testing!</li>
                <li>Design change (sidebar inspired by Arc)</li>
                <li>
                    Fixed padding issues and changed the button style (added
                    loading where necessary)
                </li>
            </ul>
            <h1>0.1.10</h1>
            <ul>
                <li>Fixed private folders issue</li>
                <li>Improved network responses to accept 403s</li>
                <li>
                    The new network response now has better structure (even
                    better soon)
                </li>
            </ul>
            <h1>0.1.9</h1>
            <ul>
                <li>Backend done for private folders (needs tidying)</li>
                <li>
                    If you visit a private folder that wasn't created by you, a
                    message will be displayed
                </li>
            </ul>
            <h1>0.1.8</h1>
            <ul>
                <li>Multiple users can be logged in at once</li>
                <li>Adding users is now <span>Admin</span> only</li>
                <li>Backend for users almost finished, UI will come later</li>
                <li>
                    New page for adding data/folders when there is no data in
                    the fetch result
                </li>
            </ul>
            <h1>0.1.7</h1>
            <ul>
                <li>Big backend changes</li>
                <li>Changed a lot of the backend to fix issues with folders</li>
                <li>Changed to a monorepo</li>
                <li>These changes are to accommodate for private folders</li>
                <li>Changed the way docCount is calculated</li>
                <li>
                    Changed JWT and cookie expiration to 10 days (may change
                    back)
                </li>
                <li>
                    Removed default docs / will add import of sample docs later
                </li>
            </ul>
            <h1>0.1.6</h1>
            <ul>
                <li>Users!</li>
                <li>Now have users which can be admin or basic</li>
                <li>
                    Admin allows all things, basic removes the ability to remove
                    things.
                </li>
                <li>Massive API overhaul (domain specific design)</li>
            </ul>
            <h1>0.1.5</h1>
            <ul>
                <li>Auth changes</li>
                <li>Now have to log in to access Blazem features</li>
                <li>
                    This is achieved by using a JWT, which is valid for 20
                    minutes (will change)
                </li>
            </ul>
            <h1>0.1.4</h1>
            <ul>
                <li>MORE API changes!</li>
                <li>Make use of Gorilla mux to complete our routing</li>
                <li>Much more efficient endpoint handling</li>
            </ul>
            <h1>0.1.3</h1>
            <ul>
                <li>API changes</li>
                <li>
                    We now do the folder filter in a better place, so if you
                    have fewer docs in a folder, it takes less time to load
                    (makes sense)
                </li>
                <li>
                    Before it would take more time if you had more docs overall
                    as we checked each doc twice.
                </li>
                <li>
                    Hopefully fixed race condition where read/write would hapeen
                    at the same time and crash.
                </li>
                <li>Fixed padding and margins</li>
            </ul>
            <h1>0.1.2</h1>
            <ul>
                <li><span>Big design change!</span></li>
                <li>Sidebar now stretches to fit the whole screen</li>
                <li>Header is now transparent to give a more modern feel</li>
            </ul>
            <h1>0.1.1</h1>
            <ul>
                <li>Added search spinner</li>
                <li>Changed panel design to fit the pattern</li>
                <li>
                    Docs in folder has two requests one for parent folders and
                    one for the data
                </li>
            </ul>
            <h1>0.1.0</h1>
            <ul>
                <li><span>Biggest API change</span></li>
                <li>Much better API, easy to add and maintain</li>
                <li>Changed sidebar design</li>
                <li>Better JSON API responses</li>
                <li>
                    Better API request checking e.g. If a POST is sent to a GET
                    only, it's rejected
                </li>
                <li>Changed design of folder and data wrappers</li>
            </ul>
            <h1>0.0.15</h1>
            <ul>
                <li>Massive design change</li>
                <li>Added better tab system</li>
                <li>
                    Added search bar for quick searches and advanced search for
                    more complex queries.
                </li>
                <li>Improved blazem backend with more efficient endpoints</li>
                <li>Added support for blazem-backup clients</li>
            </ul>
            <h1>0.0.14</h1>
            <ul>
                <li>Svelte version!</li>
                <li>Much faster load times</li>
                <li>Easier to change and improve</li>
            </ul>
            <h1>0.0.13</h1>
            <ul>
                <li>Fixed video issue where <span>.mov</span> wouldn't play</li>
                <li>changed menu item style slightly</li>
                <li>removed type view in JSON</li>
            </ul>
            <h1>0.0.12</h1>
            <ul>
                <li>Added better design for your documents</li>
                <li>Better looking on phone</li>
                <li>Started Docker version</li>
            </ul>
            <h1>0.0.11</h1>
            <ul>
                <li>Fixed nesting with where clauses</li>
                <li>Added ability to download your files.</li>
                <li>Added ability to hit 'enter' to execute queries.</li>
                <li>Changed button accents to be more easily visible.</li>
            </ul>
            <h1>0.0.10</h1>
            <ul>
                <li><span>Bug fixes:</span></li>
                <li>
                    Fixed replication bug where custom data was not replicated
                    to external cluster.
                </li>
                <li>
                    Fixed JAQL bug where nested data wasn't parsed and queried
                </li>
                <li>Fixed JAQL bug where custom JSON data wasn't queried</li>
                <li>
                    Still need to refactor JAQL code to fit on with this new
                    system.
                </li>
            </ul>
            <h1>0.0.9</h1>
            <ul>
                <li>
                    Added ability to change JSON! Can have multi-json objects,
                    this is currently dangerous as the meta-data and custom JSON
                    values are interchangeable.
                </li>
                <li>
                    Fixed replication issue where the file data wasnt
                    replicated.
                </li>
                <li>
                    Added path link so we can go back to "Your folders" as
                    suggested from user feedback.
                </li>
            </ul>
            <h1>0.0.8</h1>
            <ul>
                <li>
                    Better button design, removed yellow accent added purple
                    one.
                </li>
                <li>
                    Added 'back' functionality for when in folder (added the
                    folder path and folder href to get back to the folder page)
                </li>
            </ul>
            <h1>0.0.7</h1>
            <ul>
                <li>
                    Better data representation, shows your most important data
                    in bold (value), lists rest of JSON doc as suited.
                </li>
                <li>
                    Writes docs to disk that aren't on disk on another node. (on
                    change, this will need to change)
                </li>
            </ul>
            <h1>0.0.6</h1>
            <ul>
                <li>
                    Design change: now lists all of the JSON output of each doc
                </li>
                <li>
                    You can list specific attrs for JSON by executing SELECT
                    key,attr (have to specify key for now)
                </li>
                <li>Changed header design to look more 'tab' like.</li>
            </ul>
            <h1>0.0.5</h1>
            <ul>
                <li>
                    Files now only load once they're being viewed in modals.
                    This reduces waiting time as we're not waiting for all docs
                    to be loaded.
                </li>
                <li>
                    Folder design changed, added a placeholder image and doc
                    count.
                </li>
                <li>
                    Add Data modal improved, we now know which file is being
                    uploaded with the file name text.
                </li>
                <li>
                    Each file now has it's file name printed in the data
                    section.
                </li>
            </ul>
            <h1>0.0.4</h1>
            <ul>
                <li>
                    <span>To Disk!</span> Finally we have disk support for even further
                    data retention.
                </li>
                <li>
                    When you add a doc, it is written to disk under the <span
                        >/data</span
                    > directory.
                </li>
                <li>
                    When <span>Blazem</span> starts up, it will read the docs on
                    disk and add them to memory (this hasn't been fully tested with
                    multiple nodes yet, but is being worked on)
                </li>
                <li>
                    Changed design slighly (better menu icon, more space,
                    etc...)
                </li>
            </ul>
            <h1>0.0.3</h1>
            <ul>
                <li>
                    Redesign again: changed so we have more vertical space,
                    moved title into top bar
                </li>
                <li>
                    Added 'Run Now' button on recent queries (runs recent query)
                </li>
                <li>
                    This new design looks a lot simpler, is a lot more freindly
                    and easier to navigate.
                </li>
            </ul>
            <h1>0.0.2</h1>
            <ul>
                <li>Added a changelog :)</li>
                <li>Changed design of the Attach file in Add Data modal</li>
            </ul>
            <h1>0.0.1</h1>
            <ul>
                <li><span>First most stable version of Blazem!</span></li>
                <li>
                    Query service available to complete queries such as <span
                        >SELECT</span
                    >
                    and <span>DELETE</span>
                </li>
                <li>
                    Node <span>auto-rebalance</span> when nodes go down, recognise
                    which is master and which are followers
                </li>
                <li>Add documents to folders and edit them with modals</li>
                <li>
                    Added rules, ability to add tasks to run at specific times
                </li>
                <li>
                    Changed design to something we like, better modals, better
                    user flow
                </li>
                <li>
                    Added recent queries, ability to execute and view recently
                    executed queries
                </li>
                <li>
                    Added replication, replicate a folder to another
                    cluster/node somewhere else (is nice for testing :) )
                </li>
            </ul>
        </div>
    {/if}
</div>

<style>
    #change_log h1 {
        font-size: x-large;
    }
    #change_log h1 {
        font-size: x-large;
    }
    ul {
        list-style-type: circle;
    }
    li span {
        font-weight: 500;
    }

    #help h1 {
        font-size: x-large;
        font-weight: 500;
    }
</style>
