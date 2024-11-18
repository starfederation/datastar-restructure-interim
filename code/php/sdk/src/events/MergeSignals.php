<?php
/**
 * @copyright Copyright (c) PutYourLightsOn
 */

namespace starfederation\datastar\events;

use starfederation\datastar\Consts;
use starfederation\datastar\enums\EventType;

class MergeSignals implements EventInterface
{
    use EventTrait;

    public string $data;
    public bool $onlyIfMissing = Consts::DEFAULT_MERGE_STORE_ONLY_IF_MISSING;

    public function __construct(string $data, array $options = [])
    {
        $this->data = $data;

        foreach ($options as $key => $value) {
            $this->$key = $value;
        }
    }

    /**
     * @inerhitdoc
     */
    public function getEventType(): EventType
    {
        return EventType::MergeStore;
    }

    /**
     * @inerhitdoc
     */
    public function getDataLines(): array
    {
        $data = trim($this->data);
        $dataLines = [];

        if ($this->onlyIfMissing !== Consts::DEFAULT_MERGE_STORE_ONLY_IF_MISSING) {
            $dataLines[] = $this->getDataLine(Consts::ONLY_IF_MISSING_DATALINE_LITERAL, $this->getBooleanAsString($this->onlyIfMissing));
        }

        $lines = explode("\n", $data);
        foreach ($lines as $line) {
            $dataLines[] = $this->getDataLine(Consts::STORE_DATALINE_LITERAL, $line);
        }

        return $dataLines;
    }
}
