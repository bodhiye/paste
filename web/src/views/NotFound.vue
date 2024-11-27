<template>
    <div>
        <div class="big">
            <div class="inner">
                <img src="../assets/img/notFound.png">
                <p>{{ $t('lang.not_found.content.title') }}</p>
            </div>
        </div>
    </div>
</template>

<style scoped>
    .big{
        width: 100%;
        height: calc(100vh - 220px);
        background-color: rgb(255, 255, 255);
        position: relative;
    }
    .inner{
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
    .inner p {
        text-align: center;
        font-size: 24px;
    }
</style>
