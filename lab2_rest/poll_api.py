from fastapi.responses import JSONResponse
from fastapi import Body, FastAPI, status, exceptions
from pydantic import BaseModel
from typing import Union
import uvicorn
from fastapi import FastAPI
from enum import Enum


class VoteIn(BaseModel):
    stars: int


class VoteOut(VoteIn):
    id: int


class PollIn(BaseModel):
    name: str
    description: str | None = None


class PollOut(PollIn):
    id: int
    votes: list[VoteOut] = []


app = FastAPI()

# sample requests and queries

polls = []

poll_id = 0
vote_id = 0


def new_poll_id():
    global poll_id
    poll_id += 1
    return poll_id


def new_vote_id():
    global vote_id
    vote_id += 1
    return vote_id


@app.get("/")
async def root():
    return {"message": "Hello World"}


@app.get("/poll", response_model=list[PollOut], status_code=status.HTTP_200_OK)
async def list_polls():
    return polls


@app.post('/poll', response_model=PollOut, status_code=status.HTTP_201_CREATED)
async def create_poll(new_poll: PollIn):
    new_poll_data = new_poll.model_dump()
    new_poll_data['id'] = new_poll_id()
    new_poll_data['votes'] = []
    polls.append(new_poll_data)
    return new_poll_data


@app.get('/poll/{poll_id}', response_model=PollOut, status_code=status.HTTP_200_OK)
async def get_poll(poll_id: int):
    poll = next((p for p in polls if p['id'] == poll_id), None)
    if poll:
        return poll
    raise exceptions.HTTPException(404, 'Poll not found')


@app.put('/poll/{poll_id}', response_model=PollOut, status_code=status.HTTP_200_OK)
async def update_poll(poll_id: int, updated_poll: PollIn):
    updated_poll_data = updated_poll.model_dump()
    for i, poll in enumerate(polls):
        if poll['id'] == poll_id:
            updated_poll_data['id'] = poll_id
            updated_poll_data['votes'] = poll.get('votes', [])
            polls[i] = updated_poll_data
            return updated_poll_data

    raise exceptions.HTTPException(404, 'Poll not found')


@app.delete('/poll/{poll_id}', response_model=PollOut, status_code=status.HTTP_200_OK)
async def delete_poll(poll_id: int):
    for i, poll in enumerate(polls):
        if poll['id'] == poll_id:
            return polls.pop(i)

    raise exceptions.HTTPException(404, 'Poll not found')


@app.get('/poll/{poll_id}/vote', response_model=list[VoteOut], status_code=status.HTTP_200_OK)
async def list_poll_votes(poll_id: int):
    poll = next((p for p in polls if p['id'] == poll_id), None)
    if poll:
        return poll.get('votes', [])
    raise exceptions.HTTPException(404, 'Poll not found')


@app.post('/poll/{poll_id}/vote', response_model=VoteOut, status_code=status.HTTP_201_CREATED)
async def create_vote(poll_id: int, new_vote: VoteIn):
    new_vote_data = new_vote.model_dump()
    new_vote_data['id'] = new_vote_id()
    poll = next((p for p in polls if p['id'] == poll_id), None)
    if poll:
        poll['votes'].append(new_vote_data)
        return new_vote_data

    raise exceptions.HTTPException(404, 'Poll not found')


@app.get('/poll/{poll_id}/vote/{vote_id}', response_model=VoteOut, status_code=status.HTTP_200_OK)
async def create_vote(poll_id: int, vote_id: int):
    poll = next((p for p in polls if p['id'] == poll_id), None)
    if poll:
        vote = next((v for v in poll['votes'] if v['id'] == vote_id), None)
        if vote:
            return vote

        raise exceptions.HTTPException(404, 'Vote not found')

    raise exceptions.HTTPException(404, 'Poll not found')


@app.put('/poll/{poll_id}/vote/{vote_id}', response_model=VoteOut, status_code=status.HTTP_200_OK)
async def create_vote(poll_id: int, vote_id: int, updated_vote: VoteIn):
    updated_vote_data = updated_vote.model_dump()
    poll = next((p for p in polls if p['id'] == poll_id), None)
    if poll:
        for i, vote in enumerate(poll['votes']):
            if vote['id'] == vote_id:
                updated_vote_data['id'] = vote_id
                poll['votes'][i] = updated_vote_data
                return updated_vote_data

        raise exceptions.HTTPException(404, 'Vote not found')

    raise exceptions.HTTPException(404, 'Poll not found')


if __name__ == "__main__":
    uvicorn.run(
        "poll_api:app",
        host="0.0.0.0",
        port=8001,
        reload=True,
    )
