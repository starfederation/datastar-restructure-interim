<?php
/**
 * @copyright Copyright (c) PutYourLightsOn
 */

namespace starfederation\datastar\events;

use starfederation\datastar\Consts;
use starfederation\datastar\enums\EventType;

class ExecuteScript implements EventInterface
{
    use EventTrait;

    public string $script;
    public bool $autoRemove = Consts::DEFAULT_EXECUTE_SCRIPT_AUTO_REMOVE;
    public array $attributes = [];

    public function __construct(string $script, array $options = [])
    {
        $this->script = $script;

        foreach ($options as $key => $value) {
            $this->$key = $value;
        }
    }

    /**
     * @inerhitdoc
     */
    public function getEventType(): EventType
    {
        return EventType::ExecuteScript;
    }

    /**
     * @inerhitdoc
     */
    public function getDataLines(): array
    {
        $dataLines = [];

        if ($this->autoRemove !== Consts::DEFAULT_EXECUTE_SCRIPT_AUTO_REMOVE) {
            $dataLines[] = $this->getDataLine(Consts::AUTO_REMOVE_DATALINE_LITERAL, $this->getBooleanAsString($this->autoRemove));
        }

        foreach ($this->attributes as $attribute) {
            if ($attribute !== Consts::DEFAULT_EXECUTE_SCRIPT_ATTRIBUTES) {
                $dataLines[] = $this->getDataLine(Consts::ATTRIBUTES_DATALINE_LITERAL, $attribute);
            }
        }

        $dataLines[] = $this->getDataLine(Consts::SCRIPT_DATALINE_LITERAL, $this->script);

        return $dataLines;
    }
}
